package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	c := context.Background()
	// 1. 创建一个 3 秒超时的 context
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()

	// 等待结果或超时
	select {
	case <-ctx.Done():
		fmt.Println("Main: 任务已经超时或完成", ctx.Err())
	}
	doSomething(c)
}

func doSomething(ctx context.Context) {
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("Worker: 正在处理中...")
		case <-ctx.Done(): // 收到取消信号
			fmt.Println("Woker: 收到停止指令，优雅退出")
			return
		}
	}
}
