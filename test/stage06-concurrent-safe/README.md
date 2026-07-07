# Stage 06 测试：并发安全（sync.Mutex + 竞态检测）

还 stage05 欠下的债！把你的 CMDB API 从"单人能跑"升级到"多人并发也不崩"。
考察 `sync.Mutex` / `sync.RWMutex`，以及 `go run -race` 竞态检测工具。

**做法**：把 stage05 的 `main.go` 复制到本目录（stage06-concurrent-safe/），在它基础上改造。

## 背景场景（运维向）

你的 CMDB 上线了，10 个运维同时用。有人在查列表，有人在 POST 新增服务器。
问题来了：`net/http` 每个请求是**一个独立 goroutine**（还记得 stage03 吗？），
它们**同时读写同一个全局 `servers` 切片** —— 这就是数据竞争（data race）。

轻则数据丢失/读到脏数据，重则程序直接 panic：
`fatal error: concurrent map read and map write` 或切片扩容时崩溃。

这一关就是要：**先亲眼看到竞争，再用锁修好它。**

---

## 📚 知识铺垫（先读这个！）

### 为什么会竞争？

```go
var servers = []Server{...}   // 全局共享

// 请求 A 的 goroutine
servers = append(servers, x)  // 写

// 请求 B 的 goroutine（同一时刻）
for _, s := range servers {}  // 读
```
两个 goroutine 同时碰 `servers`，一个在写、一个在读 —— CPU 层面这不是原子操作，
读的可能读到"写了一半"的状态。这就是 data race，行为**不可预测**。

> 类比运维：两个脚本同时 `>>` 追加写同一个文件、又同时 `cat` 读它，
> 读到的内容可能是残缺的。需要 `flock` 加文件锁来串行化。
> `sync.Mutex` 就是 Go 版的 `flock`。

### sync.Mutex —— 互斥锁

```go
import "sync"

var mu sync.Mutex          // 声明一把锁（零值可用，不用 make）

mu.Lock()                  // 上锁：其他 goroutine 到这会阻塞等待
// ... 临界区：操作共享数据 ...
mu.Unlock()                // 解锁：放行下一个

// 更稳的写法：defer 保证一定解锁
mu.Lock()
defer mu.Unlock()
// ... 操作共享数据 ...
```
> 同一时刻只有一个 goroutine 能持有锁，其他排队。把并发"串行化"到临界区。
> 类比：机房只有一把钥匙，谁拿到谁进，出来才交给下一个。

### sync.RWMutex —— 读写锁（进阶）

```go
var mu sync.RWMutex

mu.RLock(); ...; mu.RUnlock()   // 读锁：多个读可以同时进行
mu.Lock();  ...; mu.Unlock()    // 写锁：独占，读写都挡在外面
```
> 适合"读多写少"场景（CMDB 就是查得多、改得少）。
> 多个"查列表"可以并发（RLock 共享），但"新增"时独占（Lock）。

### go run -race —— 竞态检测神器

```bash
go run -race main.go
```
加上 `-race`，Go 会在运行时监控内存访问，一旦发现两个 goroutine 无锁并发读写
同一变量，立刻打印 `WARNING: DATA RACE` 和出事的代码行。
> 这是你排查并发 bug 的照妖镜。stage03 提过，这关正式用上。

---

## 需求（分 3 个任务）

### 任务 1：先亲眼看到竞争（不许先加锁！）

1. 把 stage05 的 `main.go` 复制过来。
2. 写一个**并发压测脚本**（可以用 shell，也可以在 Go 里起 goroutine），
   同时发很多 GET 和 POST 请求打这个服务。
3. 用 `go run -race main.go` 启动服务，然后压测。
4. **目标**：让终端打印出 `WARNING: DATA RACE`，把它截下来/记下来。

**提示**：并发压测的简易 shell 办法（另开终端）：
```bash
for i in $(seq 1 50); do
  curl -s -X POST -d "{\"name\":\"srv-$i\",\"ip\":\"10.0.0.$i\",\"status\":\"running\"}" \
    http://localhost:8080/servers &
  curl -s http://localhost:8080/servers > /dev/null &
done
wait
```
50 个 POST 和 50 个 GET 几乎同时打，`-race` 基本必抓到竞争。

**要求**：能复现并说清楚 —— 竞争发生在哪一行？为什么？

---

### 任务 2：用 Mutex 修好它

给共享的 `servers` 加一把 `sync.Mutex`（或 `RWMutex`），
把所有读写 `servers` 的地方用锁保护起来。

**要求**：
- 定义全局锁（或把 `servers` + 锁包成一个 struct，见加分项）
- **读**（GET 列表、按 name 查）→ 加读锁 `RLock`（用 RWMutex）或 `Lock`（用 Mutex）
- **写**（POST append）→ 加写锁 `Lock`
- 每个 `Lock` 都要有对应的 `Unlock`（推荐 `defer mu.Unlock()`）
- 改完后再跑一次 `go run -race main.go` + 压测，**必须不再有 DATA RACE 警告**

**关键思考**：临界区要尽量小 —— 只锁"碰共享数据"的那几行，
别把 `json.Encode` 这种耗时操作也锁进去（会拖慢并发）。

---

### 任务 3：封装成线程安全的 Store（加分/进阶）

把"数据 + 锁"包成一个结构体，对外只暴露方法，把锁藏在内部。
这是 Go 工程里管理共享状态的标准姿势。

```go
type ServerStore struct {
    mu      sync.RWMutex
    servers []Server
}

func (s *ServerStore) List() []Server {
    s.mu.RLock()
    defer s.mu.RUnlock()
    // 返回副本，别把内部切片直接暴露出去
    return append([]Server(nil), s.servers...)
}

func (s *ServerStore) Add(srv Server) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.servers = append(s.servers, srv)
}

func (s *ServerStore) Find(name string) (Server, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, x := range s.servers {
        if x.Name == name {
            return x, true
        }
    }
    return Server{}, false
}
```

**要求**：
- handler 里不再直接碰全局切片和锁，改成调 `store.List()` / `store.Add()` / `store.Find()`
- 锁完全封装在 Store 内部，外面看不到 `mu`
- 再次 `-race` 压测通过

> 为什么这样更好？调用方不可能"忘记加锁"——锁和数据绑死在一起，
> 想访问数据只能走方法，方法内部保证了加锁。这叫"把并发安全封装进类型"。

---

## 🎯 最终要求

**通过标准：**
- 任务 1：能复现 `DATA RACE`，说清原因
- 任务 2：加锁后 `go run -race` 压测**零竞争警告**
- 任务 3（加分）：封装成 `ServerStore`，锁对外不可见
- 功能不退化：stage05 的所有接口（GET/POST/404/400/405）仍正常

**验收命令**（我检查时会跑）：
```bash
go run -race main.go        # 启动带竞态检测的服务
# 另开终端并发压测，观察有无 DATA RACE
```

---

## 💡 提示（卡住了看这里）

### 提示 1：Mutex 零值可用
```go
var mu sync.Mutex   // 直接能用，不用 make/初始化
```

### 提示 2：defer Unlock 防漏解锁
```go
mu.Lock()
defer mu.Unlock()   // 函数返回时自动解锁，即使中途 return/panic
```
不用 defer 的话，某个分支 return 忘了 Unlock → 死锁，全服务卡死。

### 提示 3：读锁能并发，写锁独占
- `RLock`：多个读 goroutine 可同时持有 → 查列表不互相阻塞
- `Lock`：写时独占 → 新增时谁也别想读写
- 选 RWMutex 更适合 CMDB（读多写少）；只用 Mutex 也对，就是读也串行了

### 提示 4：临界区要小
```go
// ❌ 锁太大，Encode 也锁进去了，拖慢并发
mu.Lock()
json.NewEncoder(w).Encode(servers)
mu.Unlock()

// ✅ 只锁"拷贝数据"这一下，Encode 在锁外做
mu.RLock()
snapshot := append([]Server(nil), servers...)
mu.RUnlock()
json.NewEncoder(w).Encode(snapshot)
```

---

写完对我说「**检查**」。我会用 `go run -race` 启动 + 并发压测，
亲眼确认你的锁真的挡住了竞争。这一关跑通，你就迈进"能写生产级并发服务"的门槛了。💪
