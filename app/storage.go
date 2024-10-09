package main

import (
	"encoding/binary"
	"fmt"

	"go.etcd.io/bbolt"
)

// Storage 封装了 bbolt.DB 操作
type Storage struct {
	db     *bbolt.DB
	bucket string
}

// NewStorage 创建并返回一个新的 Storage 实例
func NewStorage(dbFile, bucket string) (*Storage, error) {
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db, bucket: bucket}, nil
}

// EnsureBucket 确保指定的 bucket 存在
func (s *Storage) EnsureBucket() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.bucket))
		if err != nil {
			return err
		}
		return nil
	})
}

// PutString 存储一个字符串到数据库
func (s *Storage) PutString(key string, value string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		return b.Put([]byte(key), []byte(value))
	})
}

// GetString 从数据库获取一个字符串
func (s *Storage) GetString(key string) (string, error) {
	var value []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
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
func (s *Storage) PutInt(key string, value uint) error {
	intBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(intBytes, uint64(value))
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
		}
		return b.Put([]byte(key), intBytes)
	})
}

// GetInt 从数据库获取一个整数
func (s *Storage) GetInt(key string) (uint, error) {
	var intBytes []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", s.bucket)
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

// Cursor 遍历 bucket 中的所有键值对
func (s *Storage) Cursor(each func(key, value []byte)) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		if b == nil {
			return nil // 桶不存在
		}
		cursor := b.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			each(k, v) // 调用用户提供的函数指针
		}
		return nil
	})
}

// Close 关闭数据库连接
func (s *Storage) Close() error {
	return s.db.Close()
}
