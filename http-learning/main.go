package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 注册一个处理函数到路径"/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, 你好！这是我的第一个 Go HTTP 服务器！")
	})

	// 启动服务器 ， 监听8080端口
	log.Println("服务器在 http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
