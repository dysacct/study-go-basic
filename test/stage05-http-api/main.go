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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 检查请求方法（优秀的API习惯，只允许GET请求获取资源）
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"error": "method not allowed"}`))
		return
	}

	err := json.NewEncoder(w).Encode(servers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "failed to encode server list"}`))
		return
	}

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/server", serversHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("服务器查看:http://localhost:8080/server")
	fmt.Println("健康检查:http://localhost:8080/health")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
