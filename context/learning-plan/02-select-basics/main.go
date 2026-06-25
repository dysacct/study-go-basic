package main

import (
	"fmt"
	"time"
)

// 学习目标：理解 select 的多路复用机制

func main() {
	fmt.Println("=== 1. 基本的 select ===")
	basicSelect()

	fmt.Println("\n=== 2. Select 超时控制 ===")
	selectWithTimeout()

	fmt.Println("\n=== 3. Select 的 default 分支 ===")
	selectWithDefault()

	fmt.Println("\n=== 4. Select 多通道监听 ===")
	multiChannelSelect()
}

// 基本的 select：监听多个 channel
func basicSelect() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 启动两个 goroutine
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "来自 channel 1"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "来自 channel 2"
	}()

	// select 会等待第一个准备好的 channel
	select {
	case msg1 := <-ch1:
		fmt.Println("收到:", msg1)
	case msg2 := <-ch2:
		fmt.Println("收到:", msg2)
	}
}

// Select 超时控制：避免永久阻塞
func selectWithTimeout() {
	ch := make(chan string)

	// 故意不发送数据，模拟耗时操作
	go func() {
		time.Sleep(3 * time.Second) // 比超时时间长
		ch <- "这条消息不会被接收"
	}()

	select {
	case msg := <-ch:
		fmt.Println("收到:", msg)
	case <-time.After(1 * time.Second):
		fmt.Println("操作超时！")
	}
}

// Select 的 default 分支：非阻塞操作
func selectWithDefault() {
	ch := make(chan string)

	// 不启动任何 goroutine，channel 中没有数据

	select {
	case msg := <-ch:
		fmt.Println("收到:", msg)
	default:
		fmt.Println("没有数据可接收，执行默认操作")
	}
}

// Select 多通道监听：实际应用场景
func multiChannelSelect() {
	dataChannel := make(chan int)
	errorChannel := make(chan error)
	doneChannel := make(chan bool)

	// 模拟数据生产者
	go func() {
		for i := 1; i <= 3; i++ {
			time.Sleep(500 * time.Millisecond)
			dataChannel <- i
		}
		doneChannel <- true
	}()

	// 消费数据
	for {
		select {
		case data := <-dataChannel:
			fmt.Printf("处理数据: %d\n", data)
		case err := <-errorChannel:
			fmt.Println("错误:", err)
			return
		case <-doneChannel:
			fmt.Println("所有数据处理完成")
			return
		case <-time.After(2 * time.Second):
			fmt.Println("等待超时")
			return
		}
	}
}
