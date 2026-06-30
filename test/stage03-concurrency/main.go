package main

import (
	"fmt"
	"sync"
	"time"
)

type CheckResult struct {
	ServerName string
	Status     string
	Message    string
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup

	servers := []string{"web-01", "web-02", "db-01", "cache-01", "mq-01"}
	results := make(chan CheckResult, len(servers))

	for _, server := range servers {
		wg.Add(1)
		go checkServerWorker(server, 3*time.Second, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("服务器 %-8s 状态: %-7s 信息: %s\n", result.ServerName, result.Status, result.Message)
	}
	fmt.Println("Time taken:", time.Since(start))
}

func checkServer(name string, results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(2 * time.Second)

	if len(name)%2 == 0 {
		results <- CheckResult{name, "ok", "服务正常"}
	} else {
		results <- CheckResult{name, "failed", "连接超时"}
	}
	fmt.Printf("✅ %s is 检查完成\n", name)
}

func checkServerWithTimeout(name string, timeout time.Duration) CheckResult {
	// 创建一个 channel 接收检测结果
	resultCh := make(chan CheckResult, 1)

	// 在goroutine 里面执行真正的检测
	go func() {
		if name == "db-01" {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}

		if len(name)%2 == 0 {
			resultCh <- CheckResult{name, "ok", "服务正常"}
		} else {
			resultCh <- CheckResult{name, "failed", "连接超时"}
		}
	}()

	// 用 select 实现超时控制
	select {
	case result := <-resultCh:
		return result
	case <-time.After(timeout):
		return CheckResult{name, "timeout", "检测超时"}
	}
}

func checkServerWorker(name string, timeout time.Duration, results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	result := checkServerWithTimeout(name, timeout)
	results <- result
}
