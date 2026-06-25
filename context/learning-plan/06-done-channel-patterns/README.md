# Done Channel 模式：用关闭通道广播“收工”

> done channel 像店长关灯：灯一灭，所有人都知道该停手收尾。

## 本页怎么学

- 笔记文件：`context/learning-plan/06-done-channel-patterns/README.md`
- 建议节奏：先跑代码，再回来看概念；并发知识只靠读很容易“懂了，但手一写就打结”。
- 观察重点：谁发送、谁接收、谁关闭、谁负责取消、谁等待退出。

## What：核心概念

这一组笔记围绕 Go 并发控制展开：`channel` 负责 goroutine 之间传递数据或信号，`select` 负责同时等待多个通信事件，`context` 负责跨函数传播取消、超时和请求级信息，`done channel` 则是轻量的完成或停止广播。

## Why：为什么要学

真实服务里，请求可能超时，用户可能断开连接，后台任务可能需要停止。如果 goroutine 只会启动不会退出，程序就会像没关水龙头的厨房，迟早把内存和资源淹掉。

## Problem：它解决什么问题

- 避免 goroutine 泄漏。
- 让多个并发任务能协同停止。
- 给慢操作设置超时边界。
- 在调用链里传递取消信号。
- 让生产者、消费者、worker pool 等模式可控退出。

## Principle：脑内模型

- `channel` 是队列加同步点：无缓冲通道要求发送者和接收者同时到场；有缓冲通道允许先存一点数据。
- `close(ch)` 是广播“没有更多数据”，不是发送一个特殊值。
- `context.Context` 是树状传播：父 context 取消，子 context 都会收到取消。
- `done := make(chan struct{})` 常用于只传信号不传数据，关闭它可以唤醒所有等待者。
- 谁创建用于发送的通道，通常谁负责关闭；接收方不要随手关别人的通道。

## Example：关联代码

### `context/learning-plan/06-done-channel-patterns/01-basic-done.go`

```go
package main

import (
  "fmt"
  "time"
)

// 场景1：等待单个任务完成
// 实际应用：初始化操作、资源加载等

func main() {
  fmt.Println("=== 场景1：等待数据库初始化 ===")
  waitDatabaseInit()

  fmt.Println("\n=== 场景2：等待配置文件加载 ===")
  waitConfigLoad()

  fmt.Println("\n=== 场景3：等待缓存预热 ===")
  waitCacheWarmup()
}

// 等待数据库初始化完成
func waitDatabaseInit() {
  done := make(chan struct{})

  fmt.Println("开始初始化数据库...")
  go func() {
    // 模拟数据库连接和初始化
    time.Sleep(2 * time.Second)
    fmt.Println("  - 连接数据库成功")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("  - 创建连接池")
    time.Sleep(500 * time.Millisecond)
    fmt.Println("  - 执行迁移脚本")
    
    close(done) // 完成后关闭 channel
  }()

  <-done // 阻塞等待初始化完成
  fmt.Println("✓ 数据库初始化完成，应用可以启动")
}

// 等待配置文件加载
func waitConfigLoad() {
  done := make(chan struct{})
  
  fmt.Println("开始加载配置文件...")
  go func() {
    // 模拟从多个源加载配置
    configs := []string{"app.yaml", "database.yaml", "redis.yaml"}
    for _, cfg := range configs {
      time.Sleep(500 * time.Millisecond)
      fmt.Printf("  - 加载 %s\n", cfg)
    }
    
    close(done)
  }()
  
  <-done
  fmt.Println("✓ 所有配置文件加载完成")
}

// 等待缓存预热
func waitCacheWarmup() {
  done := make(chan struct{})
  
  fmt.Println("开始缓存预热...")
  go func() {
    // 模拟加载热数据到缓存
    items := []string{"用户数据", "商品列表", "分类信息", "热门文章"}
    for i, item := range items {
      time.Sleep(300 * time.Millisecond)
      fmt.Printf("  - [%d/%d] 加载 %s\n", i+1, len(items), item)
    }
    
    close(done)
  }()
  
  <-done
  fmt.Println("✓ 缓存预热完成")
}
```

### `context/learning-plan/06-done-channel-patterns/02-multiple-workers.go`

```go
package main

import (
  "fmt"
  "sync"
  "time"
)

// 场景2：等待多个 Worker 完成
// 实际应用：批量处理、并行下载、数据迁移等

func main() {
  fmt.Println("=== 场景1：批量图片处理 ===")
  batchImageProcessing()

  fmt.Println("\n=== 场景2：并发下载文件 ===")
  concurrentDownload()

  fmt.Println("\n=== 场景3：数据库批量导入 ===")
  batchDataImport()
}

// 1. 批量图片处理
func batchImageProcessing() {
  images := []string{"avatar.jpg", "banner.png", "logo.svg", "photo1.jpg", "photo2.jpg"}
  
  done := make(chan struct{})
  var wg sync.WaitGroup
  
  fmt.Printf("需要处理 %d 张图片...\n", len(images))
  
  // 启动多个 worker 处理图片
  for i, img := range images {
    wg.Add(1)
    go func(id int, filename string) {
      defer wg.Done()
      
      // 模拟图片处理：压缩、裁剪、添加水印
      time.Sleep(time.Duration(500+id*100) * time.Millisecond)
      fmt.Printf("  ✓ 处理完成: %s\n", filename)
    }(i, img)
  }
  
  // 等待所有 worker 完成后关闭 done
  go func() {
    wg.Wait()
    close(done)
  }()
  
  <-done
  fmt.Println("✓ 所有图片处理完成")
}

// 2. 并发下载文件
func concurrentDownload() {
  urls := []string{
    "https://example.com/file1.zip",
    "https://example.com/file2.pdf",
    "https://example.com/file3.mp4",
    "https://example.com/file4.doc",
  }
  
  done := make(chan struct{})
  var wg sync.WaitGroup
  
  fmt.Printf("开始下载 %d 个文件...\n", len(urls))
  
  for i, url := range urls {
    wg.Add(1)
    go downloadFile(&wg, i+1, url)
  }
  
  // 等待所有下载完成
  go func() {
    wg.Wait()
    close(done)
  }()
  
  // 主 goroutine 等待完成信号
  <-done
  fmt.Println("✓ 所有文件下载完成")
}

func downloadFile(wg *sync.WaitGroup, id int, url string) {
  defer wg.Done()
  
  fmt.Printf("  [%d] 开始下载: %s\n", id, url)
  
  // 模拟下载进度
  for progress := 0; progress <= 100; progress += 25 {
    time.Sleep(200 * time.Millisecond)
    if progress < 100 {
      fmt.Printf("  [%d] 进度: %d%%\n", id, progress)
    }
  }
  
  fmt.Printf("  [%d] ✓ 下载完成: %s\n", id, url)
}

// 3. 数据库批量导入
func batchDataImport() {
  // 模拟需要导入的数据批次
  batches := []int{1000, 2000, 1500, 3000, 2500}
  
  done := make(chan struct{})
  var wg sync.WaitGroup
  
  totalRecords := 0
  for _, count := range batches {
    totalRecords += count
  }
  
  fmt.Printf("开始导入数据，共 %d 条记录，分 %d 批...\n", totalRecords, len(batches))
  
  // 并发导入各批次
  for i, count := range batches {
    wg.Add(1)
    go importBatch(&wg, i+1, count)
  }
  
  // 等待所有批次导入完成
  go func() {
    wg.Wait()
    close(done)
  }()
  
  <-done
  fmt.Printf("✓ 数据导入完成，共导入 %d 条记录\n", totalRecords)
}

func importBatch(wg *sync.WaitGroup, batchID, count int) {
  defer wg.Done()
  
  fmt.Printf("  [批次 %d] 开始导入 %d 条记录...\n", batchID, count)
  
  // 模拟导入耗时（数据越多耗时越长）
  time.Sleep(time.Duration(count) * time.Millisecond)
  
  fmt.Printf("  [批次 %d] ✓ 导入完成\n", batchID)
}
```

### `context/learning-plan/06-done-channel-patterns/03-server-graceful-shutdown.go`

```go
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
```

### `context/learning-plan/06-done-channel-patterns/04-real-world-crawler.go`

```go
package main

import (
  "fmt"
  "math/rand"
  "sync"
  "time"
)

// 场景4：实战项目 - 并发爬虫
// 综合运用 done channel、WaitGroup、限流等技术

func main() {
  fmt.Println("=== 实战：网页爬虫系统 ===\n")
  
  crawler := NewCrawler(3) // 最多 3 个并发
  
  urls := []string{
    "https://example.com/page1",
    "https://example.com/page2",
    "https://example.com/page3",
    "https://example.com/page4",
    "https://example.com/page5",
    "https://example.com/page6",
  }
  
  results := crawler.Crawl(urls)
  
  // 统计结果
  successCount := 0
  for _, result := range results {
    if result.Success {
      successCount++
    }
  }
  
  fmt.Printf("\n=== 爬取完成 ===\n")
  fmt.Printf("总任务数: %d\n", len(urls))
  fmt.Printf("成功: %d\n", successCount)
  fmt.Printf("失败: %d\n", len(urls)-successCount)
}

// 爬虫结果
type CrawlResult struct {
  URL     string
  Success bool
  Content string
  Error   error
}

// 爬虫
type Crawler struct {
  maxConcurrent int
  semaphore     chan struct{}
}

func NewCrawler(maxConcurrent int) *Crawler {
  return &Crawler{
    maxConcurrent: maxConcurrent,
    semaphore:     make(chan struct{}, maxConcurrent),
  }
}

// 执行爬取
func (c *Crawler) Crawl(urls []string) []CrawlResult {
  done := make(chan struct{})
  resultChan := make(chan CrawlResult, len(urls))
  
  var wg sync.WaitGroup
  
  fmt.Printf("开始爬取 %d 个页面（最大并发: %d）...\n\n", len(urls), c.maxConcurrent)
  
  // 为每个 URL 启动爬取任务
  for i, url := range urls {
    wg.Add(1)
    go c.crawlURL(&wg, i+1, url, resultChan)
  }
  
  // 等待所有任务完成
  go func() {
    wg.Wait()
    close(done)
  }()
  
  // 收集结果
  results := make([]CrawlResult, 0, len(urls))
  
  // 使用 done channel 来判断是否全部完成
  go func() {
    <-done
    close(resultChan)
  }()
  
  // 收集所有结果
  for result := range resultChan {
    results = append(results, result)
  }
  
  return results
}

// 爬取单个 URL
func (c *Crawler) crawlURL(wg *sync.WaitGroup, id int, url string, resultChan chan<- CrawlResult) {
  defer wg.Done()
  
  // 获取信号量（限流）
  c.semaphore <- struct{}{}
  defer func() { <-c.semaphore }()
  
  fmt.Printf("[任务 %d] 开始爬取: %s\n", id, url)
  
  // 模拟爬取过程
  start := time.Now()
  
  // 模拟网络延迟
  delay := time.Duration(rand.Intn(2000)+500) * time.Millisecond
  time.Sleep(delay)
  
  // 模拟可能的失败（20% 失败率）
  success := rand.Float32() > 0.2
  
  result := CrawlResult{
    URL:     url,
    Success: success,
  }
  
  if success {
    result.Content = fmt.Sprintf("页面内容 (长度: %d 字节)", rand.Intn(50000)+10000)
    fmt.Printf("[任务 %d] ✓ 爬取成功: %s (耗时: %v)\n", id, url, time.Since(start))
  } else {
    result.Error = fmt.Errorf("连接超时")
    fmt.Printf("[任务 %d] ✗ 爬取失败: %s (耗时: %v)\n", id, url, time.Since(start))
  }
  
  resultChan <- result
}
```

### `context/learning-plan/06-done-channel-patterns/05-producer-consumer.go`

```go
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
// ... 后续代码请直接打开源文件继续阅读
```

## Real World Usage：真实开发场景

- HTTP 请求超时后取消数据库查询或外部 API 调用。
- 服务收到 SIGTERM 后优雅关闭 worker 和消费者。
- 批量任务用 channel 分发，用 WaitGroup 等待完成。
- 爬虫或下载器用 semaphore 限制并发量。
- 后台任务通过 done channel 停止循环。

## Common Mistakes：高频坑点

- 发送方没人接收，goroutine 永久阻塞。
- 接收方等待一个永远不会关闭的 channel。
- 多个发送者同时 close 同一个 channel 导致 panic。
- 把 context 存进结构体长期持有，而不是按调用链传递。
- 忘记调用 `cancel()`，导致计时器和资源延迟释放。
- 用 `context.Value` 传业务参数，把 context 当杂物包。

## Practice：动手练习

- 给每个 goroutine 增加退出日志，确认它真的退出。
- 把无缓冲 channel 改成有缓冲，观察阻塞点变化。
- 给耗时任务加 `context.WithTimeout`。
- 写一个 `select` 同时等待结果和超时。
- 故意不关闭 done channel，观察程序哪里卡住。

## Interview Questions：面试问答

**问：channel 关闭后还能接收吗？**

答：可以。已缓存的数据会先被读完，之后返回元素零值和 `ok=false`。

**问：context 主要解决什么问题？**

答：跨 goroutine、跨函数调用链传播取消信号、截止时间和少量请求级值。

**问：done channel 和 context 怎么选？**

答：内部简单协调可以用 done channel；跨 API 边界、需要超时或取消原因时优先用 context。

## 一句话记忆

并发代码最重要的不是“启动多少 goroutine”，而是“每个 goroutine 怎么体面地结束”。
