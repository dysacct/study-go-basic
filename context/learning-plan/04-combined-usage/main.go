package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// 学习目标：综合运用 channel、context 和 select

func main() {
	fmt.Println("=== 1. 优雅关闭多个 Worker ===")
	gracefulShutdown()

	fmt.Println("\n=== 2. 管道模式（Pipeline）===")
	pipelinePattern()

	fmt.Println("\n=== 3. 扇出扇入模式（Fan-out Fan-in）===")
	fanOutFanIn()

	fmt.Println("\n=== 4. 超时控制的 HTTP 请求模拟 ===")
	httpRequestSimulation()
}

// 1. 优雅关闭多个 Worker
func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 启动多个 worker
	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	// 等待所有 worker 退出
	<-ctx.Done()
	fmt.Println("Main: 所有 worker 已收到停止信号")
	time.Sleep(500 * time.Millisecond) // 等待打印信息
}

func worker(ctx context.Context, id int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Worker %d: 工作中...\n", id)
		case <-ctx.Done():
			fmt.Printf("Worker %d: 优雅退出\n", id)
			return
		}
	}
}

// 2. 管道模式：数据流经多个阶段处理
func pipelinePattern() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建管道：生成数字 -> 平方 -> 过滤偶数
	numbers := generate(ctx, 1, 2, 3, 4, 5)
	squares := square(ctx, numbers)
	evens := filterEven(ctx, squares)

	// 消费结果
	fmt.Print("偶数平方: ")
	for result := range evens {
		fmt.Printf("%d ", result)
	}
	fmt.Println()
}

// 生成数字
func generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// 计算平方
func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// 过滤偶数
func filterEven(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// 3. 扇出扇入模式：并行处理，聚合结果
func fanOutFanIn() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 生成任务
	jobs := make(chan int, 10)
	go func() {
		defer close(jobs)
		for i := 1; i <= 5; i++ {
			jobs <- i
		}
	}()

	// 扇出：启动多个 worker 并行处理
	numWorkers := 3
	results := make([]<-chan int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		results[i] = processJob(ctx, jobs, i+1)
	}

	// 扇入：合并所有结果
	merged := merge(ctx, results...)

	// 收集结果
	fmt.Println("处理结果:")
	for result := range merged {
		fmt.Printf("结果: %d\n", result)
	}
}

func processJob(ctx context.Context, jobs <-chan int, workerID int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for job := range jobs {
			select {
			case <-ctx.Done():
				return
			default:
				// 模拟处理时间
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
				result := job * 10
				fmt.Printf("Worker %d 处理任务 %d -> %d\n", workerID, job, result)
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

// 合并多个 channel
func merge(ctx context.Context, channels ...<-chan int) <-chan int {
	out := make(chan int)

	// 为每个输入 channel 启动一个 goroutine
	done := make(chan struct{})
	for _, ch := range channels {
		go func(c <-chan int) {
			defer func() { done <- struct{}{} }()
			for val := range c {
				select {
				case out <- val:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	// 等待所有 goroutine 完成后关闭输出
	go func() {
		for i := 0; i < len(channels); i++ {
			<-done
		}
		close(out)
	}()

	return out
}

// 4. 超时控制的 HTTP 请求模拟
func httpRequestSimulation() {
	// 设置 2 秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// 模拟 HTTP 请求
	go func() {
		// 模拟网络延迟（随机 1-3 秒）
		delay := time.Duration(rand.Intn(3)+1) * time.Second
		fmt.Printf("模拟请求，延迟 %.1f 秒...\n", delay.Seconds())

		select {
		case <-time.After(delay):
			resultCh <- "请求成功：获取到数据"
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		}
	}()

	// 等待结果或超时
	select {
	case result := <-resultCh:
		fmt.Println("✓", result)
	case err := <-errCh:
		fmt.Println("✗ 请求被取消:", err)
	case <-ctx.Done():
		fmt.Println("✗ 请求超时:", ctx.Err())
	}
}
