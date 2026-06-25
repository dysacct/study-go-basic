package main

import (
	"fmt"
	"time"
)

// 学习目标：理解 channel 的基本操作

func main() {
	fmt.Println("=== 1. 无缓冲通道 ===")
	unbufferedChannel()

	fmt.Println("\n=== 2. 有缓冲通道 ===")
	bufferedChannel()

	fmt.Println("\n=== 3. 关闭通道 ===")
	closeChannel()

	fmt.Println("\n=== 4. 单向通道 ===")
	unidirectionalChannel()
}

// 无缓冲通道：发送和接收必须同时准备好
func unbufferedChannel() {
	ch := make(chan string)

	// 启动接收者 goroutine
	go func() {
		msg := <-ch // 阻塞，直到有数据发送过来
		fmt.Println("接收到:", msg)
	}()

	// 发送数据（会阻塞，直到有接收者）
	ch <- "Hello, Channel!"
	time.Sleep(100 * time.Millisecond) // 等待 goroutine 执行完
}

// 有缓冲通道：可以在缓冲区满之前非阻塞发送
func bufferedChannel() {
	ch := make(chan int, 3) // 缓冲区大小为 3

	// 发送数据（不会阻塞，因为缓冲区未满）
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println("已发送 3 个数据，缓冲区已满")

	// 接收数据
	fmt.Println("接收:", <-ch)
	fmt.Println("接收:", <-ch)
	fmt.Println("接收:", <-ch)
}

// 关闭通道：关闭后不能再发送，但可以接收
func closeChannel() {
	ch := make(chan int, 3)

	// 发送数据
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch) // 关闭通道

	// 使用 range 遍历通道（会自动检测关闭）
	for val := range ch {
		fmt.Println("接收到:", val)
	}

	// 从已关闭的通道接收，会立即返回零值
	val, ok := <-ch
	fmt.Printf("通道已关闭: val=%d, ok=%t\n", val, ok)
}

// 单向通道：限制通道的操作方向
func unidirectionalChannel() {
	ch := make(chan string, 1)

	// 只能发送的通道
	go producer(ch)

	// 只能接收的通道
	consumer(ch)
}

// 只发送通道（chan<-）
func producer(ch chan<- string) {
	ch <- "来自生产者的消息"
	close(ch)
}

// 只接收通道（<-chan）
func consumer(ch <-chan string) {
	msg := <-ch
	fmt.Println("消费者接收:", msg)
}
