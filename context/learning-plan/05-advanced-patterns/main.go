package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// 学习目标：高级并发模式和最佳实践

func main() {
	fmt.Println("=== 1. 信号量模式：限制并发数 ===")
	semaphorePattern()

	fmt.Println("\n=== 2. 超时重试模式 ===")
	retryPattern()

	fmt.Println("\n=== 3. Context 链式传播 ===")
	contextPropagation()

	fmt.Println("\n=== 4. Done Channel 和 WaitGroup 结合 ===")
	doneWithWaitGroup()
}

// 1. 信号量模式：使用 channel 限制并发数
func semaphorePattern() {
	const maxConcurrent = 3 // 最多 3 个并发
	semaphore := make(chan struct{}, maxConcurrent)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 启动 10 个任务
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			// 获取信号量（限流）
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }() // 释放信号量
			case <-ctx.Done():
				fmt.Printf("任务 %d: 超时，未能获取执行权限\n", taskID)
				return
			}

			// 执行任务
			fmt.Printf("任务 %d: 开始执行\n", taskID)
			time.Sleep(1 * time.Second)
			fmt.Printf("任务 %d: 执行完成\n", taskID)
		}(i)
	}

	wg.Wait()
	fmt.Println("所有任务完成")
}

// 2. 超时重试模式
func retryPattern() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	maxRetries := 3
	var err error

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("尝试第 %d 次...\n", i+1)

		err = doTaskWithTimeout(ctx, 2*time.Second)
		if err == nil {
			fmt.Println("✓ 任务成功")
			return
		}

		if errors.Is(err, context.Canceled) {
			fmt.Println("✗ 整体超时，停止重试")
			return
		}

		fmt.Printf("✗ 失败: %v\n", err)

		// 等待后重试
		select {
		case <-time.After(1 * time.Second):
			continue
		case <-ctx.Done():
			fmt.Println("✗ 整体超时，停止重试")
			return
		}
	}

	fmt.Printf("✗ 达到最大重试次数: %v\n", err)
}

func doTaskWithTimeout(parentCtx context.Context, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	resultCh := make(chan error, 1)

	go func() {
		// 模拟耗时操作（有 50% 概率失败）
		time.Sleep(1500 * time.Millisecond)
		if time.Now().Unix()%2 == 0 {
			resultCh <- errors.New("操作失败")
		} else {
			resultCh <- nil
		}
	}()

	select {
	case err := <-resultCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// 3. Context 链式传播：父子关系
func contextPropagation() {
	// 根 context
	rootCtx := context.Background()

	// 第一层：设置整体超时 5 秒
	ctx1, cancel1 := context.WithTimeout(rootCtx, 5*time.Second)
	defer cancel1()

	// 第二层：添加请求 ID
	type keyType string
	ctx2 := context.WithValue(ctx1, keyType("requestID"), "req-12345")

	// 第三层：设置子任务超时 2 秒（会先于父 context 超时）
	ctx3, cancel3 := context.WithTimeout(ctx2, 2*time.Second)
	defer cancel3()

	// 执行任务
	serviceA(ctx3)
}

func serviceA(ctx context.Context) {
	type keyType string
	requestID := ctx.Value(keyType("requestID"))
	fmt.Printf("Service A 处理请求: %v\n", requestID)

	// 调用下游服务
	serviceB(ctx)
}

func serviceB(ctx context.Context) {
	type keyType string
	requestID := ctx.Value(keyType("requestID"))

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Service B 工作中... (RequestID: %v)\n", requestID)
		case <-ctx.Done():
			fmt.Printf("Service B 停止: %v\n", ctx.Err())
			return
		}
	}
}

// 4. Done Channel 和 WaitGroup 结合使用
func doneWithWaitGroup() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	done := make(chan struct{})

	// 启动多个 worker
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			workerWithContext(ctx, id)
		}(i)
	}

	// 等待所有 worker 完成
	go func() {
		wg.Wait()
		close(done) // 关闭 done channel，通知主 goroutine
	}()

	// 等待完成或超时
	select {
	case <-done:
		fmt.Println("✓ 所有 worker 正常完成")
	case <-ctx.Done():
		fmt.Println("✗ 超时，部分 worker 未完成")
		wg.Wait() // 仍然等待 worker 清理资源
	}
}

func workerWithContext(ctx context.Context, id int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < 4; i++ {
		select {
		case <-ticker.C:
			fmt.Printf("Worker %d: 步骤 %d\n", id, i+1)
		case <-ctx.Done():
			fmt.Printf("Worker %d: 提前退出\n", id)
			return
		}
	}
	fmt.Printf("Worker %d: 完成所有工作\n", id)
}
