package main

import (
	"encoding/binary"
	"fmt"
	"bytes"

	"go.etcd.io/bbolt"
)

// Storage 封装了 bbolt.DB 操作
type Storage struct {
	db     *bbolt.DB
	system *SystemStorager
	blacklist *BlacklistStorager
	counter *CounterStorager
	location *LocationStorager
	lastAccessAddr *LastAccessAddrStorager
	lastAccessIp *LastAccessIPStorager
}

// NewStorage 创建并返回一个新的 Storage 实例
func NewStorage(dbFile string) (*Storage, error) {
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	r := &Storage{db: db}
	r.EnsureBucket("system")
	r.system, err = NewSystemStorager(r, "system")
	r.EnsureBucket("blacklist")
	r.blacklist, err = NewBlacklistStorager(r, "blacklist")
	r.EnsureBucket("counter")
	r.counter, err = NewCounterStorager(r, "counter")
	r.EnsureBucket("location")
	r.location, err = NewLocationStorager(r, "location")
	r.EnsureBucket("lastAccessAddr")
	r.lastAccessAddr, err = NewLastAccessAddrStorager(r, "lastAccessAddr")
	r.EnsureBucket("lastAccessIp")
	r.lastAccessIp, err = NewLastAccessIPStorager(r, "lastAccessIp")
	return r, nil
}

// EnsureBucket 确保指定的 bucket 存在
func (s *Storage) EnsureBucket(bucket string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
}

// PutString 存储一个字符串到数据库
func (s *Storage) PutString(bucket, key , value string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		return b.Put([]byte(key), []byte(value))
	})
}

// GetString 从数据库获取一个字符串
func (s *Storage) GetString(bucket, key string) (string, error) {
	var value []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		value = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return "", err
	}
	return string(value), nil
}

// PutInt 存储一个整数到数据库
func (s *Storage) PutInt(bucket, key string, value uint) error {
	intBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(intBytes, uint64(value))
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		return b.Put([]byte(key), intBytes)
	})
}

// GetInt 从数据库获取一个整数
func (s *Storage) GetInt(bucket, key string) (uint, error) {
	var intBytes []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		intBytes = b.Get([]byte(key))
		return nil
	})
	if err != nil {
		return 0, err
	}
	if intBytes == nil {
		return 0, nil
	} else {
		return uint(binary.BigEndian.Uint64(intBytes)), nil
	}
}

// Prefix scans
func (s *Storage) PrefixScans(bucket string, prefix []byte, each func(key, value []byte)) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte(bucket)).Cursor()

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			each(k, v)
		}

		return nil
	})
}

// Range scans
func (s *Storage) RangeScans(bucket string, min, max []byte, each func(key, value []byte)) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		// Assume our events bucket exists and has RFC3339 encoded time keys.
		c := tx.Bucket([]byte(bucket)).Cursor()

		// Iterate over the 90's.
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			each(k, v)
		}

		return nil
	})
}

// 遍历 bucket 中的所有键值对
func (s *Storage) ForEach(bucket string, each func(key, value []byte)) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil // 桶不存在
		}
		b.ForEach(func(k, v []byte) error {
			each(k, v)
			return nil
		})
		return nil
	})
}

// Close 关闭数据库连接
func (s *Storage) Close() error {
	return s.db.Close()
}

type BaseStorager struct{
	delegate *Storage
	bucket string
}

func NewBaseStorager(delegate *Storage, bucket string) *BaseStorager {
	return &BaseStorager{delegate: delegate, bucket: bucket}
}

func (b BaseStorager) PutString(key string, value string) error {
	return b.delegate.PutString(b.bucket, key, value)
}
func (b BaseStorager) GetString(key string) (string, error) {
	r, err := b.delegate.GetString(b.bucket, key)
	return r, err
}
func (b BaseStorager) PutInt(key string, value uint) error {
	return b.delegate.PutInt(b.bucket, key, value)
}
func (b BaseStorager) GetInt(key string) (uint, error) {
	return b.delegate.GetInt(b.bucket, key)
}
func (b BaseStorager) PrefixScans(prefix []byte, each func(key, value []byte)) error {
	return b.delegate.PrefixScans(b.bucket, prefix, each)
}
func (b BaseStorager) RangeScans(bucket string, min, max []byte, each func(key, value []byte)) error {
	return b.delegate.RangeScans(b.bucket, min, max, each)
}
func (b BaseStorager) ForEach(each func(key, value []byte)) error {
	return b.delegate.ForEach(b.bucket, each)
}

type SystemStorager struct{
	BaseStorager
}

func NewSystemStorager(delegate *Storage, bucket string) (*SystemStorager, error) {
	return &SystemStorager{ BaseStorager: *NewBaseStorager(delegate, bucket) }, nil
}

type BlacklistStorager struct{
	BaseStorager
}

func NewBlacklistStorager(delegate *Storage, bucket string) (*BlacklistStorager, error) {
	return &BlacklistStorager{ BaseStorager: *NewBaseStorager(delegate, bucket) }, nil
}

type CounterStorager struct{
	BaseStorager
}

func NewCounterStorager(delegate *Storage, bucket string) (*CounterStorager, error) {
	return &CounterStorager{ BaseStorager: *NewBaseStorager(delegate, bucket) }, nil
}

type LocationStorager struct{
	BaseStorager
}

func NewLocationStorager(delegate *Storage, bucket string) (*LocationStorager, error) {
	return &LocationStorager{ BaseStorager: *NewBaseStorager(delegate, bucket) }, nil
}

type LastAccessAddrStorager struct{
	BaseStorager
}

func NewLastAccessAddrStorager(delegate *Storage, bucket string) (*LastAccessAddrStorager, error) {
	return &LastAccessAddrStorager{ BaseStorager: *NewBaseStorager(delegate, bucket) }, nil
}


type LastAccessIPStorager struct{
	BaseStorager
}

func NewLastAccessIPStorager(delegate *Storage, bucket string) (*LastAccessIPStorager, error) {
	return &LastAccessIPStorager{ BaseStorager: *NewBaseStorager(delegate, bucket) }, nil
}
