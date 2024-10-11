package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed dist
var dist embed.FS

// RequestHandler 接口定义
type RequestHandler interface {
	Handle(w http.ResponseWriter, r *http.Request, body io.Reader) error
}

// NewWorkConnRequestBody 结构定义
type NewWorkConnRequestBody struct {
	Content struct {
		User struct {
			User  string            `json:"user"`
			Metas map[string]string `json:"metas"`
			RunID string            `json:"run_id"`
		} `json:"user"`
		RunID        string `json:"run_id"`
		Timestamp    int64  `json:"timestamp"`
		PrivilegeKey string `json:"privilege_key"`
	} `json:"content"`
}

// NewWorkConnHandler 实现
type NewWorkConnHandler struct{}

func (h NewWorkConnHandler) Handle(w http.ResponseWriter, r *http.Request, body io.Reader) error {
	// 解析JSON请求体到一个空接口
	var newWorkConnRequest NewWorkConnRequestBody
	if err := json.NewDecoder(body).Decode(&newWorkConnRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	// fmt.Fprintf(w, "Received NewWorkConn request with content: %+v", newWorkConnRequest.Content)
	log.Printf("Received NewWorkConn request with content: %+v", newWorkConnRequest.Content)
	fmt.Fprintf(w, "{ \"reject\": false, \"unchange\": true }")
	return nil
}

// NewProxyRequestBody 结构定义
type NewProxyRequestBody struct {
	Content struct {
		User struct {
			User  string            `json:"user"`
			Metas map[string]string `json:"metas"`
			RunID string            `json:"run_id"`
		} `json:"user"`
		ProxyName          string            `json:"proxy_name"`
		ProxyType          string            `json:"proxy_type"`
		UseEncryption      bool              `json:"use_encryption"`
		UseCompression     bool              `json:"use_compression"`
		BandwidthLimit     string            `json:"bandwidth_limit"`
		BandwidthLimitMode string            `json:"bandwidth_limit_mode"`
		Group              string            `json:"group"`
		GroupKey           string            `json:"group_key"`
		RemotePort         *int              `json:"remote_port,omitempty"`    // tcp and udp only
		CustomDomains      []string          `json:"custom_domains,omitempty"` // http and https only
		Subdomain          string            `json:"subdomain,omitempty"`
		Locations          string            `json:"locations,omitempty"`
		HttpUser           string            `json:"http_user,omitempty"`
		HttpPwd            string            `json:"http_pwd,omitempty"`
		HostHeaderRewrite  string            `json:"host_header_rewrite,omitempty"`
		Headers            map[string]string `json:"headers,omitempty"`
		SK                 string            `json:"sk,omitempty"`          // stcp only
		Multiplexer        string            `json:"multiplexer,omitempty"` // tcpmux only
		Metas              map[string]string `json:"metas"`
	} `json:"content"`
}

// NewProxyHandler 实现
type NewProxyHandler struct{}

func (h NewProxyHandler) Handle(w http.ResponseWriter, r *http.Request, body io.Reader) error {
	// 解析JSON请求体到一个空接口
	var newProxyRequest NewProxyRequestBody
	if err := json.NewDecoder(body).Decode(&newProxyRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	// fmt.Fprintf(w, "Received NewProxy request with content: %+v", newProxyRequest.Content)
	log.Printf("Received NewProxy request with content: %+v", newProxyRequest.Content)
	fmt.Fprintf(w, "{ \"reject\": false, \"unchange\": true }")
	return nil
}

// NewUserConnRequestBody 结构定义
type NewUserConnRequestBody struct {
	Content struct {
		User struct {
			User  string            `json:"user"`
			Metas map[string]string `json:"metas"`
			RunID string            `json:"run_id"`
		} `json:"user"`
		ProxyName  string `json:"proxy_name"`
		ProxyType  string `json:"proxy_type"`
		RemoteAddr string `json:"remote_addr"`
	} `json:"content"`
}

// NewUserConnHandler 实现
type NewUserConnHandler struct{}

func (h NewUserConnHandler) Handle(w http.ResponseWriter, r *http.Request, body io.Reader) error {
	// 解析JSON请求体到一个空接口
	var newUserConnRequest NewUserConnRequestBody
	if err := json.NewDecoder(body).Decode(&newUserConnRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	// log.Printf("Received NewUserConn request with content: %+v", newUserConnRequest.Content)

	remoteAddr := newUserConnRequest.Content.RemoteAddr
	ip := strings.Split(remoteAddr, ":")[0]
	// 获取当前时间
	now := time.Now()
	// 记录 ip:port、访问时间
	storage.lastAccessAddr.PutString(remoteAddr, now.Format("2006-01-02 15:04:05"))
	// 记录 ip、访问时间
	storage.lastAccessIp.PutString(ip, now.Format("2006-01-02 15:04:05"))

	count, _ := storage.counter.GetInt(ip)
	storage.counter.PutInt(ip, count+1)

	fmt.Fprintf(w, "{ \"reject\": false, \"unchange\": true }")
	return nil
}

// 通用的POST请求处理函数
func PostHandler(handler RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// 限制请求体的大小
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		// 调用具体的处理逻辑
		if err := handler.Handle(w, r, r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// db 全局变量存储数据库引用
var storage *Storage

func initializeDB() error {
	// 获取用户的主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	// 定义配置目录路径
	configDir := filepath.Join(homeDir, ".config", "frp-filter")
	dbFile := "frp-filter.db"

	// 检查目录是否存在
	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		// 创建目录，包括所有必要的父目录
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		dbFile = filepath.Join(configDir, dbFile)
		log.Println("Configuration directory created:", configDir)
	} else if err != nil {
		// 处理其他错误
		log.Fatalf("Failed to stat directory: %v", err)
	} else {
		dbFile = filepath.Join(configDir, dbFile)
		log.Println("Configuration directory already exists:", configDir)
	}

	// 创建数据库包装器实例
	storage, err = NewStorage(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	password, _ := storage.system.GetString("password")
	if password == "" {
		storage.system.PutString("password", "123456")
	}

	return nil
}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	fs := http.FS(dist)

	log.Printf("Request for %s", r.URL.Path)
	// 尝试打开请求的文件
	file, err := fs.Open("dist" + r.URL.Path)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果文件不存在，检查是否请求的是根路径
			if r.URL.Path == "/" {
				// 如果请求的是根路径，尝试打开 dist/index.html
				file, err = fs.Open("dist/index.html")
				if err != nil {
					http.NotFound(w, r)
					return
				}
			} else {
				http.NotFound(w, r)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 获取文件的信息
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果请求的是目录，尝试打开目录下的 index.html
	if fileInfo.IsDir() {
		file.Close()
		file, err = fs.Open("dist" + r.URL.Path + "/index.html")
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 读取文件内容并发送给客户端
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

type AccessItem struct {
	IP    string `json:"ip"`
	Time  string `json:"time"`
	Count uint   `json:"count"`
	Info  IPInfo `json:"info"`
	// 添加更多字段...
}

func accessHandle(w http.ResponseWriter, r *http.Request) {
	// 确保只处理 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 构造要返回的数据
	data := []AccessItem{}
	storage.lastAccessIp.ForEach(func(key, value []byte) {
		ip := string(key)
		count, _ := storage.counter.GetInt(ip)
		var ipInfo *IPInfo = &IPInfo{}
		location, _ := storage.location.GetString(ip)
		if location != "" {
			json.Unmarshal([]byte(location), ipInfo)
			data = append(data, AccessItem{IP: ip, Time: string(value), Count: count, Info: *ipInfo})
		} else {
			data = append(data, AccessItem{IP: ip, Time: string(value), Count: count})
		}
	})

	// 设置响应的内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 将数据编码为 JSON 并写入响应
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ipLocationHandle(w http.ResponseWriter, r *http.Request) {
	// 确保只处理 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	query := r.URL.Query()

	// 获取单个参数值
	ip := query.Get("ip")

	location, _ := storage.location.GetString(ip)
	if location == "" {
		iPInfoResponse, _ := GetIPInfo(ip)
		if iPInfoResponse.Data == nil {
			location = "{}"
		} else {
			data, _ := json.Marshal(iPInfoResponse.Data)
			location = string(data)
			storage.location.PutString(ip, location)
		}
	}

	// 设置响应的内容类型为 JSON
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(location))
}

func initializeHandle() {
	// 注册处理器
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/access", accessHandle)
	http.HandleFunc("/ip/location", ipLocationHandle)
	http.Handle("/new_proxy", PostHandler(NewProxyHandler{}))
	http.Handle("/new_work_conn", PostHandler(NewWorkConnHandler{}))
	http.Handle("/new_user_conn", PostHandler(NewUserConnHandler{}))

}

func initialize() error {
	if err := initializeDB(); err != nil {
		return err
	}
	initializeHandle()
	return nil
}

func main() {
	port := flag.String("port", "8000", "define the port to listen on")
	flag.Parse() // 解析命令行参数

	if err := initialize(); err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
		return
	}

	defer func() {
		if err := storage.Close(); err != nil {
			log.Fatalf("Failed to close database: %v", err)
		}
	}()

	// 使用提供的端口号启动服务器
	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Starting server on %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
