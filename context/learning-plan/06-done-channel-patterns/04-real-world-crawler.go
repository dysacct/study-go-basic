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

