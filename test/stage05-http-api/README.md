# Stage 05 测试：HTTP API 服务（net/http + encoding/json）

这一关对着你的核心目标项目 **HTTP API 服务 + 简易 CMDB** 去。
考察标准库 `net/http`、`encoding/json`，以及前面学的 struct、map、错误处理。

在本目录新建 `main.go`，写一个**服务器资产管理 API**（简易 CMDB 雏形）。

## 背景场景（运维向）

你想管理一批服务器信息（名字、IP、状态），不想每次都 SSH 上去看。
干脆做个 HTTP API：`curl` 一下就能查所有服务器、查单台、加一台。
这就是 CMDB（配置管理数据库）最朴素的样子，也是所有后端服务的基本盘。

---

## 📚 知识铺垫（先读这个！）

### net/http 最小服务器
```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // 注册路由：路径 -> 处理函数
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "ok")   // 往 w 写就是返回给客户端的内容
    })

    fmt.Println("服务启动在 :8080")
    http.ListenAndServe(":8080", nil)  // 阻塞监听
}
```
跑起来后：`curl http://localhost:8080/health` → 返回 `ok`

> 类比运维：`http.HandleFunc` 就像 nginx 的 `location` 块——把某个路径映射到某个处理逻辑。
> `w http.ResponseWriter` = 响应（你往里写 = 返回给客户端），`r *http.Request` = 请求（里面有方法、路径、参数、body）。

### 返回 JSON
```go
type Server struct {
    Name   string `json:"name"`     // 反引号里是 tag，控制 JSON 字段名
    IP     string `json:"ip"`
    Status string `json:"status"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    s := Server{Name: "web-01", IP: "10.0.0.5", Status: "running"}
    w.Header().Set("Content-Type", "application/json")  // 告诉客户端这是 JSON
    json.NewEncoder(w).Encode(s)                        // 把结构体编码成 JSON 写给客户端
}
```

### 读取请求信息
```go
r.Method                          // 请求方法："GET" / "POST" ...
r.URL.Query().Get("name")         // 查询参数：?name=web-01 取出 "web-01"
json.NewDecoder(r.Body).Decode(&s) // 把请求 body 的 JSON 解析进结构体
```

### JSON 编解码对照
| 方向 | 函数 | 类比 |
|------|------|------|
| 结构体 → JSON | `json.Marshal` / `Encoder.Encode` | 序列化，发出去 |
| JSON → 结构体 | `json.Unmarshal` / `Decoder.Decode` | 反序列化，收进来 |

---

## 需求（分 4 个任务，逐步加功能）

### 任务 1：跑起一个 HTTP 服务

实现一个 `/health` 接口，`curl` 访问返回 `ok`（纯文本即可）。

**要求**：
- 用 `http.HandleFunc` 注册 `/health`
- 用 `http.ListenAndServe(":8080", nil)` 启动
- 启动时在终端打印一句提示
- 验证：`curl http://localhost:8080/health` 返回 `ok`

---

### 任务 2：返回 JSON —— 查所有服务器

准备一份"内存数据"（用切片模拟数据库），实现 `/servers` 返回**所有服务器的 JSON 列表**。

```go
type Server struct {
    Name   string `json:"name"`
    IP     string `json:"ip"`
    Status string `json:"status"`
}

// 内存里的假数据（真实项目里是数据库）
var servers = []Server{
    {"web-01", "10.0.0.5", "running"},
    {"db-01", "10.0.0.21", "running"},
    {"cache-01", "10.0.0.30", "stopped"},
}
```

**要求**：
- `/servers` 返回整个切片的 JSON 数组
- 设置 `Content-Type: application/json`
- 用 struct tag 让 JSON 字段是小写的 `name`/`ip`/`status`
- 验证：`curl http://localhost:8080/servers` 返回 JSON 数组

---

### 任务 3：查询参数 —— 按名字查单台

实现 `/server?name=web-01`，根据查询参数返回**匹配的那一台**服务器。

**要求**：
- 用 `r.URL.Query().Get("name")` 取参数
- 在 `servers` 切片里查找匹配的
- 找到 → 返回那台的 JSON
- 没找到 → 返回 HTTP 404 状态码 + 错误提示
  提示：`http.Error(w, "server not found", http.StatusNotFound)`
- 参数为空 → 返回 400（提示要带 name 参数）
- 验证：
  - `curl "http://localhost:8080/server?name=web-01"` → 返回 web-01
  - `curl "http://localhost:8080/server?name=xxx"` → 404

---

### 任务 4：POST 新增服务器（区分请求方法）

实现 `/servers` 的 **POST** 方法：接收 JSON body，把新服务器加进内存切片。
（GET `/servers` 保持任务 2 的查询功能，POST 用来新增——同一路径，按方法区分。）

**要求**：
- 用 `r.Method` 判断：`GET` 走查询，`POST` 走新增，其他方法返回 405
- POST 时用 `json.NewDecoder(r.Body).Decode(&newServer)` 解析 body
- 解析失败（body 不是合法 JSON）→ 返回 400
- 成功 → `append` 进 `servers`，返回新增的那台 JSON（状态码 201 Created）
- 验证：
  ```bash
  curl -X POST http://localhost:8080/servers \
    -H "Content-Type: application/json" \
    -d '{"name":"mq-01","ip":"10.0.0.40","status":"running"}'
  ```
  然后再 `curl http://localhost:8080/servers` 应该能看到多出来的 mq-01

---

## 🎯 最终要求

`main.go` 实现上面 4 个接口，一个服务同时提供：
- `GET /health` → ok
- `GET /servers` → 所有服务器
- `GET /server?name=xxx` → 单台
- `POST /servers` → 新增

**通过标准：**
- 服务能启动，各接口 `curl` 能通
- JSON 编解码正确（字段名小写）
- 404 / 400 / 405 等状态码用对
- GET/POST 按方法正确分流
- 代码整洁（变量用驼峰命名，别再下划线了 😄）

---

## 💡 提示（卡住了看这里）

### 提示 1：怎么测试？
服务启动后**会一直阻塞**（这是正常的，Web 服务本来就常驻）。
**另开一个终端**用 `curl` 测，或者浏览器直接访问 GET 接口。
停服务：`Ctrl+C`。

### 提示 2：按方法分流的骨架
```go
http.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // 返回列表
    case http.MethodPost:
        // 解析 body，新增
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
})
```

### 提示 3：状态码常量（别背数字，用常量）
```go
http.StatusOK                  // 200
http.StatusCreated             // 201
http.StatusBadRequest          // 400
http.StatusNotFound            // 404
http.StatusMethodNotAllowed    // 405
```

### 提示 4：设置状态码的顺序（重要！）
```go
w.Header().Set("Content-Type", "application/json")  // ① 先设 header
w.WriteHeader(http.StatusCreated)                   // ② 再写状态码
json.NewEncoder(w).Encode(s)                        // ③ 最后写 body
```
顺序反了状态码会不生效（这是 net/http 的一个经典坑）。

---

写完对我说「**检查**」。检查时我会**启动你的服务并用 curl 实测**每个接口。这一关跑通，你就有一个真能用的 API 服务了，成就感拉满！💪
