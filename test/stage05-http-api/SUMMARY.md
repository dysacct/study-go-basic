# Stage 05 总结：HTTP API 服务（net/http + encoding/json）

## 这一关做了什么
写了一个**服务器资产管理 API**（简易 CMDB 雏形）：健康检查、查列表、查单台、POST 新增，完整的 HTTP 状态码体系。这是目标项目"HTTP API 服务 + 简易 CMDB"的骨架。

## 核心知识点

### 1. 最小 HTTP 服务
```go
mux := http.NewServeMux()              // 路由器（推荐用它，别用默认的 nil）
mux.HandleFunc("/health", handler)     // 注册：路径 -> 处理函数
http.ListenAndServe(":8080", mux)      // 阻塞监听
```
> `http.HandleFunc` 类比 nginx 的 `location` 块。
> `w http.ResponseWriter` = 响应（往里写=返回给客户端），`r *http.Request` = 请求。
> ⚠️ 一个路径只能注册一次，重复注册 → panic。

### 2. 返回 JSON
```go
type Server struct {
    Name string `json:"name"`   // 反引号 tag 控制 JSON 字段名（小写）
}
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(servers)   // 结构体/切片 -> JSON 写给客户端
```

### 3. 读请求
```go
r.Method                              // "GET"/"POST"...
r.URL.Query().Get("name")             // ?name=xxx 取参数
json.NewDecoder(r.Body).Decode(&obj)  // 请求 body 的 JSON -> 结构体
```

### 4. 按方法分流（RESTful 骨架）
```go
switch r.Method {
case http.MethodGet:    // 查询
case http.MethodPost:   // 新增
default:
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
```

### 5. 状态码用常量
```go
http.StatusOK(200) StatusCreated(201) StatusBadRequest(400)
http.StatusNotFound(404) StatusMethodNotAllowed(405)
```

## 关键领悟（血泪教训）

### 领悟 1：WriteHeader 是"一锤子买卖"，顺序错了状态码失效 ⭐
```go
// ❌ 错误：还没检查方法就把状态码焊死成 200
w.WriteHeader(http.StatusOK)
if r.Method != "GET" {
    w.WriteHeader(405)   // 无效！Go 报 "superfluous WriteHeader call"，状态码仍是 200
}
```
**规矩：先做完所有检查、确定要返回啥，最后才 WriteHeader。**
黄金顺序：**设 header → 写状态码 → 写 body**。
> 类比：快递面单贴上传送带就发走了，之后改不了。面单必须贴之前写对。

### 领悟 2：http.Error 顺序天生正确
`http.Error(w, msg, code)` 内部一步到位（设 header + 写状态码 + 写 body），
返回错误时优先用它，不会踩顺序坑。

### 领悟 3：json.Encoder.Encode 自动补 200
返回正常列表时不用手动 `WriteHeader(200)`，Encode 会自动发。
只有要返回非 200（如 201 Created）时才手动 WriteHeader。

### 领悟 4：每个请求是独立 goroutine → 共享数据要加锁 ⭐
```go
servers = append(servers, newServer)   // 多个请求并发 POST = 数据竞争！
```
`net/http` 每个请求开一个 goroutine（呼应 stage03）。多 goroutine 同时写全局切片
会 data race。生产级要加 `sync.Mutex` 或用 channel 串行化。
用 `go run -race main.go` 可检测。

## 加分技巧：业务字段校验
JSON 格式合法 ≠ 数据合法。解析成功后还要查必填字段：
```go
if newServer.Name == "" || newServer.IP == "" {
    http.Error(w, "missing required fields", http.StatusBadRequest)
}
```

## 测试方法
- 服务启动后**常驻阻塞**（正常）。**另开终端**用 curl 测。
- POST 测试：`curl -X POST -d '{"name":"mq-01",...}' http://localhost:8080/servers`
- 看状态码：`curl -w "[HTTP %{http_code}]\n" ...`

## 一句话记忆
- 一个路径只注册一次，用 `http.NewServeMux()`
- 状态码顺序：**header → WriteHeader → body**，且 WriteHeader 只一次
- 返回错误用 `http.Error`（顺序天生对），返回列表用 `Encode`（自动 200）
- 共享数据 + 并发请求 = 记得加锁（stage03 的债，迟早要还）
