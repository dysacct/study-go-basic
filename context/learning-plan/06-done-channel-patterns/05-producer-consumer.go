package main

import (
  "fmt"
  "sync"
  "time"
)

// 场景5：生产者-消费者模式
// 实际应用：日志处理、消息队列、任务调度等

func main() {
  fmt.Println("=== 场景1：日志处理系统 ===")
  logProcessingSystem()

  fmt.Println("\n=== 场景2：订单处理系统 ===")
  orderProcessingSystem()
}

// 1. 日志处理系统
func logProcessingSystem() {
  logChan := make(chan string, 10)
  done := make(chan struct{})
  
  var wg sync.WaitGroup
  
  // 启动 3 个日志生产者（不同的服务模块）
  modules := []string{"API", "Database", "Cache"}
  for _, module := range modules {
    wg.Add(1)
    go logProducer(&wg, logChan, module, done)
  }
  
  // 启动 2 个日志消费者（写入文件、发送到监控系统）
  consumers := []string{"FileWriter", "MonitorSender"}
  for _, name := range consumers {
    wg.Add(1)
    go logConsumer(&wg, logChan, name, done)
  }
  
  fmt.Println("日志系统已启动\n")
  
  // 运行 5 秒后关闭
  time.Sleep(5 * time.Second)
  fmt.Println("\n正在关闭日志系统...")
  close(done)
  
  // 等待生产者停止
  time.Sleep(100 * time.Millisecond)
  close(logChan) // 关闭日志通道
  
  // 等待消费者处理完剩余日志
  wg.Wait()
  fmt.Println("✓ 日志系统已关闭")
}

func logProducer(wg *sync.WaitGroup, logChan chan<- string, module string, done <-chan struct{}) {
  defer wg.Done()
  
  ticker := time.NewTicker(800 * time.Millisecond)
  defer ticker.Stop()
  
  logID := 0
  
  for {
    select {
    case <-ticker.C:
      logID++
      log := fmt.Sprintf("[%s] 日志消息 #%d", module, logID)
      
      select {
      case logChan <- log:
        // 发送成功
      default:
        fmt.Printf("[%s] ⚠️  日志队列已满，丢弃日志\n", module)
      }
      
    case <-done:
      fmt.Printf("[%s Producer] 停止生产日志\n", module)
      return
    }
  }
}

func logConsumer(wg *sync.WaitGroup, logChan <-chan string, name string, done <-chan struct{}) {
  defer wg.Done()
  
  processedCount := 0
  
  for {
    select {
    case log, ok := <-logChan:
      if !ok {
        fmt.Printf("[%s] 日志通道已关闭，共处理 %d 条日志\n", name, processedCount)
        return
      }
      
      // 模拟处理日志
      time.Sleep(100 * time.Millisecond)
      processedCount++
      fmt.Printf("[%s] 处理: %s\n", name, log)
      
    case <-done:
      // 收到关闭信号，但继续处理剩余日志
      fmt.Printf("[%s] 收到关闭信号，处理剩余日志...\n", name)
    }
  }
}

// 2. 订单处理系统
func orderProcessingSystem() {
  orderChan := make(chan Order, 20)
  done := make(chan struct{})
  
  var wg sync.WaitGroup
  
  // 启动订单生产者（模拟用户下单）
  wg.Add(1)
  go orderProducer(&wg, orderChan, done)
  
  // 启动多个订单处理器
  workerCount := 3
  for i := 1; i <= workerCount; i++ {
    wg.Add(1)
    go orderProcessor(&wg, i, orderChan, done)
  }
  
  fmt.Println("订单处理系统已启动\n")
  
  // 运行 6 秒后关闭
  time.Sleep(6 * time.Second)
  fmt.Println("\n停止接收新订单...")
  close(done)
  
  // 等待生产者停止
  time.Sleep(100 * time.Millisecond)
  close(orderChan)
  
  // 等待所有订单处理完成
  fmt.Println("处理剩余订单...")
  wg.Wait()
  fmt.Println("✓ 所有订单处理完成，系统已关闭")
}

type Order struct {
  ID       int
  UserID   string
  Amount   float64
  CreateAt time.Time
}

func orderProducer(wg *sync.WaitGroup, orderChan chan<- Order, done <-chan struct{}) {
  defer wg.Done()
  
  orderID := 1000
  ticker := time.NewTicker(700 * time.Millisecond)
  defer ticker.Stop()
  
  for {
    select {
    case <-ticker.C:
      order := Order{
        ID:       orderID,
        UserID:   fmt.Sprintf("user%d", orderID%100),
        Amount:   float64((orderID%1000 + 100)),
        CreateAt: time.Now(),
      }
      
      select {
      case orderChan <- order:
        fmt.Printf("📝 新订单: #%d (金额: ¥%.2f)\n", order.ID, order.Amount)
        orderID++
      default:
        fmt.Println("⚠️  订单队列已满，请稍后重试")
      }
      
    case <-done:
      fmt.Println("[Producer] 停止接收新订单")
      return
    }
  }
}

func orderProcessor(wg *sync.WaitGroup, workerID int, orderChan <-chan Order, done <-chan struct{}) {
  defer wg.Done()
  
  processedCount := 0
  totalAmount := 0.0
  
  for {
    select {
    case order, ok := <-orderChan:
      if !ok {
        fmt.Printf("[Worker %d] 完成，共处理 %d 个订单，总金额: ¥%.2f\n", 
          workerID, processedCount, totalAmount)
        return
      }
      
      // 模拟订单处理（验证、扣款、发货等）
      fmt.Printf("[Worker %d] 处理订单 #%d...\n", workerID, order.ID)
      time.Sleep(time.Duration(500+workerID*100) * time.Millisecond)
      
      processedCount++
      totalAmount += order.Amount
      fmt.Printf("[Worker %d] ✓ 订单 #%d 处理完成\n", workerID, order.ID)
      
    case <-done:
      // 继续处理剩余订单
    }
  }
}

