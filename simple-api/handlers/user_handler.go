package handlers

import (
	"encoding/json"
	"net/http"
	"simple-api/data"
	"strconv"
	"strings"
)

// UserHandler 处理用户相关的路由
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头为 JSON
	w.Header().Set("Content-Type", "application/json")

	// 只允许 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// 去除路径开头的 /users，得到剩余部分
	path := strings.TrimPrefix(r.URL.Path, "/users")
	path = strings.TrimSuffix(path, "/")

	// 如果路径为空或只有 "/"，表示请求所有用户
	if path == "" {
		users := data.GetAllUsers()
		json.NewEncoder(w).Encode(users)
		return
	}

	// 否则尝试解析ID
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
	}

	user := data.GetUserByID(id)
	if user == nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(user)
}
