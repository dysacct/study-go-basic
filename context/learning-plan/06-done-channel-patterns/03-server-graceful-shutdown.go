package main

import (
  "fmt"
  "os"
  "os/signal"
  "sync"
  "syscall"
  "time"
)

// 场景3：服务优雅关闭
// 实际应用：Web 服务器、消息队列消费者、定时任务等

func main() {
  fmt.Println("=== 模拟 Web 服务器优雅关闭 ===\n")
  runWebServer()
}

// 模拟一个完整的 Web 服务器生命周期
func runWebServer() {
  // 创建关闭信号 channel
  done := make(chan struct{})
  
  // 监听系统信号
  sigChan := make(chan os.Signal, 1)
  signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
  
  // 启动各种服务组件
  var wg sync.WaitGroup
  
  // 1. HTTP 服务器
  wg.Add(1)
  go runHTTPServer(&wg, done)
  
  // 2. 定时任务
  wg.Add(1)
  go runScheduledTasks(&wg, done)
  
  // 3. 消息队列消费者
  wg.Add(1)
  go runMessageConsumer(&wg, done)
  
  // 4. 后台任务处理器
  wg.Add(1)
  go runBackgroundWorker(&wg, done)
  
  fmt.Println("✓ 所有服务组件已启动")
  fmt.Println("📝 提示：按 Ctrl+C 触发优雅关闭\n")
  
  // 等待关闭信号（或者 5 秒后自动触发，方便演示）
  select {
  case <-sigChan:
    fmt.Println("\n⚠️  收到关闭信号，开始优雅关闭...")
  case <-time.After(5 * time.Second):
    fmt.Println("\n⚠️  演示时间到，触发优雅关闭...")
  }
  
  // 关闭 done channel，通知所有组件停止
  close(done)
  
  // 等待所有组件完成清理工作
  fmt.Println("等待所有服务组件完成清理...")
  wg.Wait()
  
  fmt.Println("\n✓ 所有服务已安全关闭，程序退出")
}

// HTTP 服务器
func runHTTPServer(wg *sync.WaitGroup, done <-chan struct{}) {
  defer wg.Done()
  
  fmt.Println("[HTTP Server] 启动，监听 :8080")
  
  ticker := time.NewTicker(1 * time.Second)
  defer ticker.Stop()
  
  requestCount := 0
  
  for {
    select {
    case <-ticker.C:
      requestCount++
      fmt.Printf("[HTTP Server] 处理请求 #%d\n", requestCount)
      
    case <-done:
      fmt.Println("[HTTP Server] 收到关闭信号")
      fmt.Println("[HTTP Server] 等待现有请求完成...")
      time.Sleep(500 * time.Millisecond) // 模拟处理中的请求
      fmt.Println("[HTTP Server] ✓ 已关闭")
      return
    }
  }
}

// 定时任务
func runScheduledTasks(wg *sync.WaitGroup, done <-chan struct{}) {
  defer wg.Done()
  
  fmt.Println("[Scheduler] 启动，每 2 秒执行清理任务")
  
  ticker := time.NewTicker(2 * time.Second)
  defer ticker.Stop()
  
  taskCount := 0
  
  for {
    select {
    case <-ticker.C:
      taskCount++
      fmt.Printf("[Scheduler] 执行清理任务 #%d\n", taskCount)
      
    case <-done:
      fmt.Println("[Scheduler] 收到关闭信号")
      fmt.Println("[Scheduler] 完成当前任务后退出...")
      time.Sleep(300 * time.Millisecond)
      fmt.Println("[Scheduler] ✓ 已关闭")
      return
    }
  }
}

// 消息队列消费者
func runMessageConsumer(wg *sync.WaitGroup, done <-chan struct{}) {
  defer wg.Done()
  
  fmt.Println("[MQ Consumer] 启动，消费消息队列")
  
  ticker := time.NewTicker(1500 * time.Millisecond)
  defer ticker.Stop()
  
  msgCount := 0
  
  for {
    select {
    case <-ticker.C:
      msgCount++
      fmt.Printf("[MQ Consumer] 消费消息 #%d\n", msgCount)
      
    case <-done:
      fmt.Println("[MQ Consumer] 收到关闭信号")
      fmt.Println("[MQ Consumer] 处理缓冲区剩余消息...")
      time.Sleep(400 * time.Millisecond)
      fmt.Println("[MQ Consumer] 提交消费偏移量...")
      time.Sleep(200 * time.Millisecond)
      fmt.Println("[MQ Consumer] ✓ 已关闭")
      return
    }
  }
}

// 后台任务处理器
func runBackgroundWorker(wg *sync.WaitGroup, done <-chan struct{}) {
  defer wg.Done()
  
  fmt.Println("[Worker] 启动，处理后台任务")
  
  ticker := time.NewTicker(2500 * time.Millisecond)
  defer ticker.Stop()
  
  jobCount := 0
  
  for {
    select {
    case <-ticker.C:
      jobCount++
      fmt.Printf("[Worker] 处理任务 #%d\n", jobCount)
      
    case <-done:
      fmt.Println("[Worker] 收到关闭信号")
      fmt.Println("[Worker] 完成当前任务...")
      time.Sleep(600 * time.Millisecond)
      fmt.Println("[Worker] 保存工作进度...")
      time.Sleep(200 * time.Millisecond)
      fmt.Println("[Worker] ✓ 已关闭")
      return
    }
  }
}

