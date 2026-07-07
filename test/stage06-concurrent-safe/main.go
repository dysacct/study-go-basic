package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var mu sync.RWMutex

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

// GET /servers
// POST /servers
func serversHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		// -------- 临界区开始（只复制共享数据）--------
		mu.RLock()
		serverList := append([]Server(nil), servers...)
		mu.RUnlock()
		// -------- 临界区结束 --------

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(serverList); err != nil {
			http.Error(w, "failed to encode server list", http.StatusInternalServerError)
			return
		}

	case http.MethodPost:

		var newServer Server

		if err := json.NewDecoder(r.Body).Decode(&newServer); err != nil {
			http.Error(w, "Bad Request: Invalid JSON Body", http.StatusBadRequest)
			return
		}

		if newServer.Name == "" ||
			newServer.IP == "" ||
			newServer.Status == "" {

			http.Error(w, "Bad Request: Missing required fields", http.StatusBadRequest)
			return
		}

		// -------- 临界区开始（写共享数据）--------
		mu.Lock()
		servers = append(servers, newServer)
		mu.Unlock()
		// -------- 临界区结束 --------

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_ = json.NewEncoder(w).Encode(newServer)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// GET /server?name=xxx
func serverNameHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "missing required parameter 'name'", http.StatusBadRequest)
		return
	}

	var (
		result Server
		found  bool
	)

	// -------- 临界区开始（只读共享数据）--------
	mu.RLock()

	for _, s := range servers {
		if s.Name == name {
			result = s // 拷贝出来
			found = true
			break
		}
	}

	mu.RUnlock()
	// -------- 临界区结束 --------

	if !found {
		http.Error(w, "server not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/servers", serversHandler)
	mux.HandleFunc("/server", serverNameHandler)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	fmt.Println("服务器启动：http://localhost:8080/servers")
	fmt.Println("按名称查询：http://localhost:8080/server?name=web-01")
	fmt.Println("健康检查：http://localhost:8080/health")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
