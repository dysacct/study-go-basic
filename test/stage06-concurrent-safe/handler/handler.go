package handler

import (
	"encoding/json"
	"net/http"
	"project-001/model"
	"project-001/store"
)

type Handler struct {
	Store *store.ServerStore
}

// 这一步是把 store.ServerStore 传进来，Handler 就可以直接调用 Store 的方法了。
func New(store *store.ServerStore) *Handler {
	return &Handler{
		Store: store,
	}
}

// 这一步是把 Handler 的方法注册到路由上，Handler 就可以直接处理请求了。
func (h *Handler) Servers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := h.Store.List()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	case http.MethodPost:
		var server model.Server

		if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// 业务字段校验：JSON 合法 ≠ 数据合法，必填字段不能空
		if server.Name == "" || server.IP == "" || server.Status == "" {
			http.Error(w, "missing required fields", http.StatusBadRequest)
			return
		}

		h.Store.Add(server)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(server)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// 这一步是把 Handler 的方法注册到路由上，Handler 就可以直接处理请求了。
func (h *Handler) Server(w http.ResponseWriter, r *http.Request) {

	// 只允许 GET 查询
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")

	// 参数为空：是"没传参数"（客户端的错），返回 400，
	// 而不是当成"没找到这台"（404）——两种语义要分清
	if name == "" {
		http.Error(w, "missing required parameter 'name'", http.StatusBadRequest)
		return
	}

	server, ok := h.Store.Find(name)

	if !ok {
		http.Error(w, "server not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(server)
}
