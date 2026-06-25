package main

import (
	"context"
	"fmt"
	"time"
)

// 学习目标：理解 context 的四种创建方式和使用场景

func main() {
	fmt.Println("=== 1. WithCancel: 手动取消 ===")
	withCancel()

	fmt.Println("\n=== 2. WithTimeout: 超时自动取消 ===")
	withTimeout()

	fmt.Println("\n=== 3. WithDeadline: 指定截止时间 ===")
	withDeadline()

	fmt.Println("\n=== 4. WithValue: 传递请求范围的数据 ===")
	withValue()
}

// WithCancel: 可以手动取消的 context
func withCancel() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine: 收到取消信号，退出")
				return
			default:
				fmt.Println("Goroutine: 工作中...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 2 秒后手动取消
	time.Sleep(2 * time.Second)
	fmt.Println("Main: 发送取消信号")
	cancel()
	time.Sleep(1 * time.Second) // 等待 goroutine 退出
}

// WithTimeout: 超时自动取消
func withTimeout() {
	// 创建一个 2 秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 建议总是调用 cancel 释放资源

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine: 超时退出, 原因:", ctx.Err())
				return
			default:
				fmt.Println("Goroutine: 处理中...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	// 等待超时
	<-ctx.Done()
	fmt.Println("Main: 检测到超时")
	time.Sleep(500 * time.Millisecond)
}

// WithDeadline: 指定截止时间
func withDeadline() {
	// 设置 3 秒后的截止时间
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Goroutine: 到达截止时间，退出")
				return
			case <-ticker.C:
				// 检查剩余时间
				if deadline, ok := ctx.Deadline(); ok {
					remaining := time.Until(deadline)
					fmt.Printf("Goroutine: 剩余时间 %.1f 秒\n", remaining.Seconds())
				}
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("Main: 截止时间已到")
	time.Sleep(500 * time.Millisecond)
}

// WithValue: 传递请求范围的数据
func withValue() {
	// 创建一个携带值的 context
	type keyType string
	const userIDKey keyType = "userID"
	const traceIDKey keyType = "traceID"

	ctx := context.WithValue(context.Background(), userIDKey, "user123")
	ctx = context.WithValue(ctx, traceIDKey, "trace-abc-456")

	// 传递给处理函数
	processRequest(ctx, userIDKey, traceIDKey)
}

func processRequest(ctx context.Context, userKey, traceKey interface{}) {
	// 从 context 中获取值
	userID := ctx.Value(userKey)
	traceID := ctx.Value(traceKey)

	fmt.Printf("处理请求 - UserID: %v, TraceID: %v\n", userID, traceID)

	// 可以继续传递给其他函数
	doSomework(ctx, userKey)
}

func doSomework(ctx context.Context, userKey interface{}) {
	userID := ctx.Value(userKey)
	fmt.Printf("执行工作 - UserID: %v\n", userID)
}
