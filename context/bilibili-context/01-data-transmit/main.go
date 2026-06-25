package main

import (
	"context"
	"fmt"
)

type UserInfo struct {
	Name string
}

func main() {
	c := context.Background() // == context.TODO()
	c = context.WithValue(c, "abc", UserInfo{Name: "白茶"})
	GetUser(c)
}

func GetUser(ctx context.Context) {
	fmt.Println(ctx.Value("abc").(UserInfo).Name)
}
