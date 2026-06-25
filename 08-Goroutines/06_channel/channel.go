package main

import (
	"fmt"
	"time"
)

func prepareDrink(orderChannel chan string) {
	fmt.Println("Barista: Starting to prepare drink...")
	time.Sleep(3 * time.Second)
	fmt.Println("Barista: Done!")
	orderChannel <- "Latte" // 将字符串Latte写入通道
}

func main() {
	orderChannel := make(chan string) // 创建一个通道内存空间
	//var orderChannel chan string

	fmt.Println("Coffee shop opens...")
	go prepareDrink(orderChannel)

	order := <-orderChannel // waiting here till Latte is ready
	fmt.Println("Order ready!", order)
	fmt.Println("Coffee shop closes...")
}
