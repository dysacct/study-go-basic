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

