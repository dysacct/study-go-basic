# Stage 07 测试：HTTP 客户端 + 并发采集（Prometheus 指标雏形）

前六关你写的都是**服务端**（收请求、回响应）。这一关掉个头——
你要写**客户端**：主动去"抓"别人的 HTTP 接口，并发抓一堆目标，再把结果
按 Prometheus 格式暴露出去。这就是 **Prometheus 采集器 / blackbox exporter** 的雏形，
是你目标项目「Prometheus 指标采集」的正式起点。

考察：`net/http` 客户端、`http.Client` 超时、`io.ReadAll`、并发采集
（复用 stage03 的 goroutine/WaitGroup + stage06 的并发安全）、
Prometheus 文本格式解析（复用 stage04 的 bufio/strings）。

在本目录新建 `main.go`。

## 背景场景（运维向）

你管着一批服务，想知道它们**活没活着、响应多快**。总不能一台台手动 `curl`。
你要写个采集器：给它一个 URL 列表，它**并发**去探测每一个，记录
"通不通、HTTP 状态码、响应耗时"，最后把这些数据变成 Prometheus 能抓的指标。
Prometheus 定时来抓你的 `/metrics`，Grafana 画图，告警系统盯着——监控闭环的第一环就是你。

> 类比：这就是 Go 版的 `for h in hosts; do curl -w '%{time_total}' $h & done`，
> 但带超时控制、错误隔离、结果聚合，而且并发安全。

---

## 📚 知识铺垫（先读这个！）

### 1. HTTP 客户端 —— 主动发请求

服务端是 `http.ListenAndServe`（等别人来）；客户端是 `http.Get`（主动去发）。

```go
import (
    "io"
    "net/http"
    "time"
)

// ❌ 别用 http.Get() 裸奔——它没有超时，目标卡住你就永远挂着
// ✅ 自己建带超时的 Client（生产铁律：任何网络请求都要有超时）
client := &http.Client{
    Timeout: 3 * time.Second,   // 整个请求 3 秒不完成就放弃
}

resp, err := client.Get("http://localhost:8080/health")
if err != nil {
    // 连不上、超时、DNS 失败等——网络层错误在这里
    return
}
defer resp.Body.Close()          // ⚠️ 必须关！不然连接泄漏（像文件不 Close）

body, err := io.ReadAll(resp.Body)  // 把响应体全部读出来（[]byte）
if err != nil {
    return
}

resp.StatusCode                  // HTTP 状态码（200/404/500...）int 类型
string(body)                     // 响应体转字符串
```

> **三个必记的坑：**
> 1. `http.Get()` 无超时，别用；用 `http.Client{Timeout}`。
> 2. `resp.Body` 是个流，**必须 `defer resp.Body.Close()`**，否则连接泄漏。
> 3. `err == nil` 只代表"网络通了"，不代表成功——**HTTP 500 也是 err==nil**，
>    要另外看 `resp.StatusCode`。（类比：电话打通了 ≠ 对方答应了你）

### 2. 并发采集 —— 复用 stage03 + stage06

50 个目标串行抓，每个 3 秒 = 150 秒。并发抓 = 最慢那个的时间。

```go
var wg sync.WaitGroup
var mu sync.Mutex               // 保护共享的 results（stage06 的教训！）
results := make([]Result, 0)

for _, url := range targets {
    wg.Add(1)
    go func(u string) {         // ⚠️ 把 url 当参数传进去（闭包变量陷阱，见提示4）
        defer wg.Done()
        r := probe(client, u)   // 探测单个目标
        mu.Lock()
        results = append(results, r)   // 并发写切片 → 必须加锁
        mu.Unlock()
    }(url)
}
wg.Wait()                       // 等所有 goroutine 抓完
```

> **错误隔离**：一个目标挂了（超时/拒连），不能影响其他目标的采集。
> 每个 goroutine 各自处理自己的 error，把"失败"也当成一种结果记下来（success=false）。

### 3. Prometheus 文本格式 —— 采集器的"通用语"

Prometheus 的指标就是**纯文本**，一行一个，格式极简：

```
# HELP probe_success 目标是否探测成功
# TYPE probe_success gauge
probe_success{target="http://localhost:8080/health"} 1
probe_success{target="http://10.0.0.99/health"} 0
probe_duration_seconds{target="http://localhost:8080/health"} 0.023
```

规则：
- `#` 开头是注释（HELP 说明、TYPE 类型），可选但推荐
- 数据行格式：`指标名{标签1="值",标签2="值"} 数值`
- `{...}` 里是标签（label），用来区分同名指标的不同维度（哪个 target）
- 最后是数值（1=成功/0=失败，或耗时秒数）

> 解析它，正是 stage04 的活：`bufio.Scanner` 逐行读，`strings` 切割，
> 跳过 `#` 和空行。生成它，就是 `fmt.Fprintf` 按格式拼字符串。

---

## 需求（分 4 个任务，逐步搭出一个采集器）

### 任务 1：HTTP 客户端 —— 探测单个目标

写一个函数，用带超时的 `http.Client` 去 GET 一个 URL，返回探测结果。

```go
type Result struct {
    Target     string        // 目标 URL
    Success    bool          // 是否成功（能连通且状态码 2xx）
    StatusCode int           // HTTP 状态码（失败时为 0）
    Duration   time.Duration // 响应耗时
    Err        string        // 错误信息（成功时为空）
}

func probe(client *http.Client, url string) Result {
    // 1. 记录开始时间 time.Now()
    // 2. client.Get(url)，defer resp.Body.Close()
    // 3. 出错 → Success=false，记下 Err
    // 4. 成功 → 记 StatusCode、Duration，判断 2xx 算 Success
}
```

**要求**：
- 用 `http.Client{Timeout: 3*time.Second}`，不许用裸 `http.Get()`
- 记得 `defer resp.Body.Close()`
- `main` 里先硬编码探测一个目标（如你 stage06 的 `http://localhost:8080/health`），打印结果
- 验证：先启动 stage06 的服务，再跑本程序，应打印出 success/状态码/耗时

---

### 任务 2：并发采集多个目标（错误隔离 + 并发安全）

给一个 URL 列表，**并发**探测所有目标，聚合结果。

```go
var targets = []string{
    "http://localhost:8080/health",   // 正常（stage06 起着）
    "http://localhost:8080/servers",  // 正常
    "http://localhost:9999/nope",     // 故意连不上，测错误隔离
}
```

**要求**：
- 用 `goroutine` + `sync.WaitGroup` 并发抓（复用 stage03）
- 结果写进共享切片，用 `sync.Mutex` 保护（复用 stage06——**别再裸 append 了**）
- **错误隔离**：9999 那个连不上，不能影响另外两个正常返回
- 把 goroutine 里的 `url` 用**参数**传进去（别踩闭包变量陷阱，见提示 4）
- 全部抓完后，打印每个目标的结果（成功的/失败的都要显示）
- 验证：`go run -race main.go`（并发写切片，用 -race 确认你锁对了！）

---

### 任务 3：解析 Prometheus 文本格式

真实场景里，你抓到的目标（比如 node_exporter）返回的就是 Prometheus 文本。
写个函数把它解析出来，能查某个指标的值。

给你一段样例文本（可以硬编码成字符串，或放个 `.txt` 文件读）：
```
# HELP node_cpu_seconds_total CPU time
# TYPE node_cpu_seconds_total counter
node_cpu_seconds_total{mode="idle"} 12345.6
node_cpu_seconds_total{mode="user"} 678.9
node_memory_free_bytes 8589934592
```

**要求**：
- 用 `bufio.Scanner` 逐行读（复用 stage04）
- **跳过** `#` 开头的注释行和空行
- 解析出每行的"指标名"和"数值"（标签部分能拆出来更好，拆不出至少要能拿到数值）
- 提供一个查询：给指标名，返回它的数值（如查 `node_memory_free_bytes` → 8589934592）
- 用 `strconv.ParseFloat` 把字符串数值转成 `float64`
- 验证：解析上面的文本，打印所有指标名+值，并能查出指定指标

---

### 任务 4：做成 Exporter —— 暴露自己的 /metrics（综合）

把前面串起来：起一个 HTTP 服务，`/metrics` 端点**实时采集**所有目标，
把结果按 **Prometheus 文本格式**吐出来。这就是一个能被 Prometheus 抓取的 exporter。

**要求**：
- 复用 stage05 的 HTTP 服务骨架，注册 `/metrics`
- 每次访问 `/metrics` 时，并发采集 `targets`（复用任务 2）
- 把结果格式化成 Prometheus 文本，比如：
  ```
  # HELP probe_success 目标探测是否成功 (1=成功, 0=失败)
  # TYPE probe_success gauge
  probe_success{target="http://localhost:8080/health"} 1
  probe_success{target="http://localhost:9999/nope"} 0
  # HELP probe_duration_seconds 探测耗时(秒)
  # TYPE probe_duration_seconds gauge
  probe_duration_seconds{target="http://localhost:8080/health"} 0.0231
  ```
- `Content-Type` 设为 `text/plain`（Prometheus 文本格式的标准）
- 耗时用秒（float），`fmt.Fprintf(w, "...%.4f\n", d.Seconds())`
- 验证：`curl http://localhost:8081/metrics`（注意换个端口，别和 stage06 的 8080 撞）
  应看到标准 Prometheus 格式的指标输出

---

## 🎯 最终要求

`main.go` 实现一个采集器：
- `probe()` 探测单目标（带超时、错误处理）
- 并发采集多目标（WaitGroup + Mutex，-race 通过）
- 能解析 Prometheus 文本格式
- `/metrics` 端点按 Prometheus 格式暴露采集结果

**通过标准：**
- 探测正常目标返回 success，探测挂掉的目标返回 failure（错误隔离生效）
- `go run -race` 并发采集无 DATA RACE
- `/metrics` 输出合法的 Prometheus 文本
- 任何网络请求都有超时，`resp.Body` 都有 `defer Close()`

**验收命令**（我检查时会跑）：
```bash
# 先起 stage06 的服务当采集目标
cd ../stage06-concurrent-safe && go run . &
# 再起你的采集器
go run -race main.go
curl http://localhost:8081/metrics
```

---

## 💡 提示（卡住了看这里）

### 提示 1：客户端 vs 服务端，别搞混
| | 服务端（前几关） | 客户端（这关） |
|---|---|---|
| 角色 | 等别人来（被动） | 主动去发（主动） |
| 核心 | `http.ListenAndServe` | `client.Get(url)` |
| 数据 | `w http.ResponseWriter` 往里写 | `resp.Body` 往外读 |

### 提示 2：err==nil 不代表成功
```go
resp, err := client.Get(url)
// err==nil 只说明"网络通了、拿到响应了"
// HTTP 404/500 时 err 依然是 nil！要单独判断：
if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
    success = true
}
```

### 提示 3：Body 一定要 Close，且顺序有讲究
```go
resp, err := client.Get(url)
if err != nil {
    return   // ← err!=nil 时 resp 可能是 nil，不能 defer Close，直接返回
}
defer resp.Body.Close()   // ← 确认 err==nil 后再 defer
```

### 提示 4：goroutine 闭包变量陷阱（经典坑）⭐
```go
// ❌ 错误：所有 goroutine 共享同一个 url 变量，循环跑完 url 已是最后一个
for _, url := range targets {
    go func() { probe(client, url) }()   // 大概率都抓最后一个目标！
}

// ✅ 正确：把 url 当参数传进去，每个 goroutine 拿到自己的副本
for _, url := range targets {
    go func(u string) { probe(client, u) }(url)
}
```
> 注：Go 1.22+ 修复了 for 循环变量语义，新版本可能不踩这坑。但**显式传参**是
> 永远安全的好习惯，别赌 Go 版本。

### 提示 5：解析一行 Prometheus 文本
```go
line := `node_memory_free_bytes 8589934592`
// 最朴素：按空格切成"指标部分"和"数值部分"
// 更细：用 strings.Fields，或找最后一个空格 strings.LastIndex(line, " ")
// 有标签的行 node_cpu{mode="idle"} 123：指标名在 { 前，数值在 } 后的空格后
```

### 提示 6：格式化输出 Prometheus 文本
```go
fmt.Fprintf(w, "# HELP probe_success 目标探测是否成功\n")
fmt.Fprintf(w, "# TYPE probe_success gauge\n")
for _, r := range results {
    val := 0
    if r.Success { val = 1 }
    fmt.Fprintf(w, "probe_success{target=%q} %d\n", r.Target, val)
}
```
> `%q` 会自动给字符串加双引号,正好符合 label 的格式 `target="..."`。

---

写完对我说「**检查**」。我会启动你的采集器 + 一个目标服务，
用 `-race` 验证并发安全，`curl /metrics` 看 Prometheus 格式对不对。
这一关跑通，你就有了监控体系最核心的那个"探针"，离真正的 Prometheus 集成只差一层窗户纸了。💪
