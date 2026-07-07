package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Status string `json:"status"`
}

var servers = []Server{
	{"web-01", "10.0.0.5", "running"},
	{"db-01", "10.0.0.21", "running"},
	{"cache-01", "10.0.0.30", "stopped"},
}

func serversHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//// 检查请求方法（优秀的API习惯，只允许GET请求获取资源）
		//if r.Method != http.MethodGet {
		//	w.WriteHeader(http.StatusMethodNotAllowed)
		//	_, _ = w.Write([]byte(`{"error": "method not allowed"}`))
		//	return
		//}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(servers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error": "failed to encode server list"}`))
			return
		}
	case http.MethodPost:
		var newServer Server

		// 使用流式解码器直接从请求体 r.Body 中解析 JSON
		// 传入指针 &newServer 以便将解析后的值直接写入该结构体变量
		err := json.NewDecoder(r.Body).Decode(&newServer)
		if err != nil {
			http.Error(w, "Bad Request: Invalid JSON Body", http.StatusBadRequest)
			return
		}
		if newServer.Name == "" || newServer.IP == "" || newServer.Status == "" {
			http.Error(w, "Bad Request: Missing required fields", http.StatusBadRequest)
			return
		}

		servers = append(servers, newServer)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(newServer)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func serverNameHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 安全校验： 只允许 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. 提取并清洗查询参数
	// r.URL.Query() 会把 URL 里的参数解析成一个 map[string][]string
	// .Get("name") 会精准去除第一个匹配的值，如果参数不存在，返回空字符串 ""
	name := r.URL.Query().Get("name")

	// 3. 边界校验: 参数为空则直接返回 HTTP 400 错误
	if name == "" {
		http.Error(w, "bad request: missing required parameter 'name'", http.StatusBadRequest)
		return
	}
	// 4. 遍历切片查找匹配的服务器
	for _, s := range servers {
		if s.Name == name {
			// 找到匹配资产, 设置响应头并返回 JSON
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(s)
			return // 找到后立即拦截并返回, 终止函数
		}
	}

	http.Error(w, "server not found", http.StatusNotFound)
}

//func serversNameHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
//		return
//	}
//	w.Header().Set("Content-Type", "application/json")
//	_ = json.NewEncoder(w).Encode(servers)
//}

func main() {
	mux := http.NewServeMux()
	//mux.HandleFunc("/servers", serversHandler)
	mux.HandleFunc("/server", serverNameHandler)
	mux.HandleFunc("/servers", serversHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("服务器启动:curl http://localhost:8080/servers")
	fmt.Println("服务器启动:http://localhost:8080/server?name=web-01")
	fmt.Println("健康检查:http://localhost:8080/health")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
