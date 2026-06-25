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

