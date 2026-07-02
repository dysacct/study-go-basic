# Stage 03 总结：并发（goroutine / channel / select / context）

## 这一关做了什么
写了一个**批量服务器健康检测工具**：并发检测多台服务器，收集结果，加超时控制，加整体取消。这是你最薄弱的一块，也是最有含金量的一关。

## 并发四件套

| 工具 | 作用 | 类比 Shell |
|------|------|-----------|
| **goroutine** | 开启并发任务 | `command &` 后台执行 |
| **channel** | goroutine 间传数据 | 命名管道 `mkfifo` |
| **WaitGroup** | 等所有任务完成 | `wait` |
| **select** | 多路复用（超时/取消） | `select` 等多个 fd |
| **context** | 统一控制取消/超时 | `kill -TERM -PID`（进程组）|

## 核心知识点

### 1. goroutine + WaitGroup
```go
var wg sync.WaitGroup
for _, s := range servers {
    wg.Add(1)               // 计数 +1（在 go 之前！）
    go func() {
        defer wg.Done()     // 完成时 -1
        // 干活
    }()
}
wg.Wait()                   // 阻塞到计数归零
```

### 2. channel 收集结果 + 关闭时机
```go
results := make(chan CheckResult, len(servers))  // 带缓冲，防死锁
go func() {
    wg.Wait()          // 等所有 worker 干完
    close(results)     // 再关闭 channel（关键！）
}()
for r := range results {   // 读到 close 才结束
    ...
}
```
> `for range channel` 必须配 `close`，否则读空后永久阻塞（deadlock）。

### 3. select + time.After —— 单个任务超时
```go
select {
case r := <-resultCh:            // 结果先到
    return r
case <-time.After(3*time.Second): // 超时先到
    return timeoutResult
}
```

### 4. context —— 整体取消
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()                   // 必须！释放定时器资源

select {
case r := <-resultCh:
    return r
case <-ctx.Done():               // 整体超时/取消，一键通知所有人
    return cancelledResult
}
```

## 关键领悟（面试常问）

### 领悟 1：缓冲 channel 防 goroutine 泄漏
```go
resultCh := make(chan CheckResult, 1)  // 缓冲 1 是精髓
```
超时后主流程走了，但内部 goroutine 还在跑。缓冲 1 让它醒来后能把结果塞进缓冲区、正常退出，**不会永久阻塞变僵尸**。无缓冲（0）就会泄漏。

### 领悟 2：defer cancel() 为什么必须
`WithTimeout` 内部起了个定时器（后台计时进程）。任务提前干完时，`cancel()` 把它 kill 掉回收资源。不调用 → 泄漏。`go vet` 会专门警告。
> 类比：`trap 'cleanup' EXIT`——保证无论怎么退出都清理。

### 领悟 3：time.After vs ctx.Done() 的区别
| | time.After | ctx.Done() |
|--|-----------|-----------|
| 粒度 | 每个任务各自一个闹钟 | 一个信号通知所有人 |
| 类比 | 每个任务单独 `timeout 3s` | `kill -TERM -进程组` |
| 场景 | 单台超时 | Ctrl+C / 全局截止 |

## 踩过的坑
- 一开始"写出来了但逻辑不熟"——并发就是这样，**写对≠想通**，靠画数据流图逐帧理解。
- 死函数没删干净（注释掉旧版本时漏了个 `}`）。

## 一句话记忆
- `wg.Add` 在 `go` **之前**；`close` 在 `wg.Wait` **之后**
- 超时模式的 channel 用**缓冲 1** 防泄漏
- `WithTimeout` 后立刻 **`defer cancel()`**
