# Stage 03 测试：并发编程（goroutine + channel + context）

考察你 `08-Goroutines` 和 `context` 目录学的并发知识。  
在本目录新建 `main.go`，实现一个**批量服务器健康检测工具**（运维刚需）。

## 背景场景（运维向）

你管着 10 台服务器，要批量检测它们的健康状态。  
串行检测太慢（每台 2 秒，总共 20 秒），必须用**并发**。  
而且要能**超时控制**（某台卡住了不能拖累全局）、**收集结果**（最后统一汇报）。

---

## 📚 知识铺垫（先读这个！）

### 并发三件套的分工

| 工具 | 作用 | 类比 Shell |
|------|------|------------|
| **goroutine** | 开启并发任务 | `command &` 后台执行 |
| **channel** | goroutine 间传递数据 | 命名管道 `mkfifo` |
| **WaitGroup** | 等所有任务完成 | `wait` 等所有后台任务 |
| **select** | 多路复用（超时、取消） | `select` 等多个 fd |
| **context** | 统一控制取消/超时 | `kill -TERM -PID`（进程组） |

### 常见坑（运维人最容易踩的）

1. **忘记 `wg.Wait()`** → main 退出了，goroutine 还没跑完就被杀了  
   就像 Shell 脚本里 `command &` 后没 `wait`，进程直接退出

2. **channel 忘记关闭** → `for range ch` 会永久阻塞  
   就像管道写端没关，读端 `read` 一直等

3. **往已关闭的 channel 写入** → panic  
   就像往已删除的管道文件 `echo`

4. **goroutine 泄漏** → 起了 goroutine 但没法退出，越积越多  
   就像后台进程没清理，`ps aux` 一堆僵尸进程

---

## 需求（分 4 个任务，逐步加功能）

### 任务 1：基础并发（goroutine + WaitGroup）

定义一个函数模拟检测单台服务器：
```go
func checkServer(name string, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // 模拟耗时操作（真实场景是 SSH 连接或 HTTP 请求）
    time.Sleep(2 * time.Second)
    
    fmt.Printf("✅ %s 检测完成\n", name)
}
```

在 `main` 里并发检测 5 台服务器：
```go
servers := []string{"web-01", "web-02", "db-01", "cache-01", "mq-01"}
```

**要求**：
- 用 `sync.WaitGroup` 等所有 goroutine 完成
- 打印开始时间和结束时间，证明是并发的（总耗时约 2 秒，不是 10

---

### 任务 2：收集结果（加上 channel）

改造 `checkServer`，返回检测结果到 channel：
```go
type CheckResult struct {
    ServerName string
    Status     string  // "ok" 或 "failed"
    Message    string
}

func checkServer(name string, results chan<- CheckResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
    time.Sleep(2 * time.Second)
    
    // 模拟：随机失败（用服务器名长度判断，偶数=成功，奇数=失败）
    if len(name) % 2 == 0 {
        results <- CheckResult{name, "ok", "服务正常"}
    } else {
        results <- CheckResult{name, "failed", "连接超时"}
    }
}
```

在 `main` 里：
1. 创建 channel：`results := make(chan CheckResult, len(servers))`
2. 启动所有 goroutine
3. 等待完成后**关闭 channel**（重点！）
4. 用 `for range results` 读取所有结果并打印

**要求**：
- channel 要带缓冲（`make(chan CheckResult, N)`），否则可能死锁
- 必须在 `wg.Wait()` 之后关闭 channel
- 打印出每台服务器的状态

---

### 任务 3：单个任务超时控制（select + time.After）

有的服务器可能卡住很久，不能让它拖累全局。给单个检测加 **3 秒超时**：

```go
func checkServerWithTimeout(name string, timeout time.Duration) CheckResult {
    // 创建一个 channel 接收检测结果
    resultCh := make(chan CheckResult, 1)
    
    // 在 goroutine 里执行真正的检测
    go func() {
        time.Sleep(2 * time.Second)  // 模拟检测耗时
        
        if len(name) % 2 == 0 {
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
```

在 `main` 里调用这个新函数，证明超时机制有效。  
**提示**：可以把某台服务器的 `time.Sleep` 改成 5 秒，触发超时。

---

### 任务 4：整体取消控制（context）

用户按 Ctrl+C 或者整体超时了，要能**立即取消所有正在进行的检测**。

改造成带 context 的版本：
```go
func checkServerWithContext(ctx context.Context, name string) CheckResult {
    resultCh := make(chan CheckResult, 1)
    
    go func() {
        time.Sleep(2 * time.Second)
        
        if len(name) % 2 == 0 {
            resultCh <- CheckResult{name, "ok", "服务正常"}
        } else {
            resultCh <- CheckResult{name, "failed", "连接超时"}
        }
    }()
    
    select {
    case result := <-resultCh:
        return result
    case <-ctx.Done():
        return CheckResult{name, "cancelled", "检测被取消"}
    }
}
```

在 `main` 里：
1. 创建一个 5 秒超时的 context：`ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)`
2. 别忘了 `defer cancel()`
3. 把 ctx 传给所有检测函数
4. 如果某几台检测很慢（比如 10 秒），5 秒后应该被 context 取消

---

## 🎯 最终要求

`main.go` 里实现上面 4 个任务（可以分成 4 个函数 `task1()`, `task2()`, `task3()`, `task4()` 依次调用，也可以只保留最终版 task4）。

**必须能跑通并输出：**
- 并发执行的证据（总耗时远小于串行）
- 所有服务器的检测结果
- 超时/取消机制生效的证据

**通过标准：**
- 没有 goroutine 泄漏（检测完能正常退出）
- 没有 deadlock（channel 使用正确）
- context 取消能立即生效
- 代码整洁，注释清晰

---

## 💡 提示（卡住了看这里）

### 提示 1：WaitGroup 的正确姿势
```go
var wg sync.WaitGroup
for _, name := range servers {
    wg.Add(1)
    go checkServer(name, &wg)  // 传指针
}
wg.Wait()  // 阻塞到所有 Done 完成
```

### 提示 2：channel 关闭时机
```go
go func() {
    wg.Wait()      // 等所有 goroutine 完成
    close(results) // 然后关闭 channel
}()

for result := range results {  // 读取直到 channel 关闭
    fmt.Println(result)
}
```

### 提示 3：select 语法
```go
select {
case msg := <-ch:
    // 收到消息
case <-time.After(3 * time.Second):
    // 超时了
case <-ctx.Done():
    // 被取消了
}
```

---

写完对我说「**检查**」。如果卡住了**随时问我**，并发这块一开始绕很正常，我会手把手带你过。加油！💪
