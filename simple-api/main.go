package main

import (
	"log"
	"net/http"
	"simple-api/handlers"
)

func main() {
	// 注册路由: 所有以/users 开头的请求都交给 UserHandler 处理
	http.HandleFunc("/users", handlers.UserHandler)
	http.HandleFunc("/users/", handlers.UserHandler)

	// 可选：添加一个根路径的欢迎信息
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(`{"message": "Welcome to Simple API! Visit /users to see data."}`))
			return
		}
		http.NotFound(w, r)
	})

	// 启动服务器
	log.Println("Server starting on http://localhost:8080")
	log.Println("Try: GET http://localhost:8080/users")
	log.Println("Try: POST http://localhost:8080/users/1")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
