package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 定义测试密码
	password := map[string]string{
		"admin": "abcd123",
		"test":  "123456",
	}

	fmt.Println("=== 生成密码 Hash (bcrypt) ===")

	// 遍历每一个账号, 生成Hash
	for username, password := range password {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10) // 10 是难度系数, 越高越安全, 但越慢，最高40
		if err != nil {
			log.Fatalf("生成失败: %v", err)
		}
		fmt.Printf("用户: %s\n", username)
		fmt.Printf("密码: %s\n", password)
		fmt.Printf("Hash: %s\n", string(hash))
		fmt.Println("--------------------------------")
	}

	fmt.Println("=== 生成密码 Hash (bcrypt) 完成 ===")
}
