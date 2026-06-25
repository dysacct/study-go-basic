package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	user := User{Name: "Bogdan", Age: 30}
	// 序列化
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data)) // string 方法是[]byte 转换为 string
}
