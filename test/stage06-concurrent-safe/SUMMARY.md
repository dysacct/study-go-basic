# Stage 06 总结：并发安全（sync.RWMutex + 竞态检测 + 封装）

## 这一关做了什么
把 stage05 的 CMDB API 从"单人能跑"升级到"多人并发也不崩"：拆成 model/store/handler/main 四包分层，用 `sync.RWMutex` 把共享的 `servers` 保护起来，并封装成线程安全的 `ServerStore`。用 `go run -race` 亲手验证零竞争。这是"能写生产级并发服务"的门槛。

## 核心知识点

### 1. 每个 HTTP 请求 = 一个 goroutine（地基）
`net/http` 每来一个请求就单独开一个 goroutine 处理。所以 N 个请求 = N 个 goroutine **同时**碰共享数据 → 不加锁必 data race。（呼应 stage03）

### 2. 数据竞争的本质：append 不是原子操作
```go
servers = append(servers, x)  // CPU 眼里分 3 步：读len → 写槽位 → 改len
```
两个 goroutine 交错执行这 3 步 → 数据被覆盖丢失；扩容时搬家 → 读到野内存 → panic。

### 3. sync.RWMutex —— 读写锁（读多写少首选）
```go
mu.RLock(); ...; mu.RUnlock()   // 读锁：多个读可并发
mu.Lock();  ...; mu.Unlock()    // 写锁：独占，读写全挡外面
```
规则：读+读放行 ✅；读+写、写+写互斥 🚫。CMDB 查多改少，用 RWMutex 比 Mutex 并发好。

### 4. defer Unlock —— 防死锁
```go
mu.Lock()
defer mu.Unlock()   // 无论从哪个出口离开（return/panic），锁一定释放
```
忘了解锁 → 后续 goroutine 全卡在 Lock() → 整个服务死锁。

### 5. go run -race —— 竞态照妖镜
```bash
go run -race .   // 运行时监视内存访问，发现无锁并发读写立刻 WARNING: DATA RACE + 指出行号
```
并发 bug 时有时无，不能靠手测，必须用 -race 逼它现形。

## 关键领悟（血泪教训）

### 领悟 1：含锁结构体必须传指针 ⭐
```go
func New() *ServerStore { ... }   // 返回 *ServerStore，不是 ServerStore
```
按值传会**连锁一起复制**成两把独立的锁 → 锁失效，等于没锁。`go vet` 会自动报"锁被复制"。

### 领悟 2：List() 要返回副本，不能返回内部切片 ⭐
```go
return append([]model.Server(nil), s.servers...)   // 拷一份快照返回
```
直接 `return s.servers` 等于把内部数组地址交出去；锁释放后调用方还在读，内部 append 扩容搬家 → 脏读/race。返回副本让"锁的保护窗口"和"数据生命周期"对齐。
> 类比：同事问有哪些机器，别扔钥匙让他进机房自己数（你还在插拔），要拍张快照给他。

### 领悟 3：临界区要尽量小 ⭐
```go
list := store.List()              // 锁内只做拷贝（纳秒级）
json.NewEncoder(w).Encode(list)   // 慢活（序列化+写socket）在锁外做
```
别把 json.Encode 这种耗时 IO 锁进去，否则持锁太久，其他读全堵住。

### 领悟 4：把并发安全封装进类型（最高级）⭐
```go
type ServerStore struct {
    mu      sync.RWMutex     // 小写 = 私有，包外看不见
    servers []model.Server   // 小写 = 私有，包外碰不到
}
```
锁和数据都私有，外面**只能**走 List/Add/Find 方法，方法内部已加好锁 → **调用方想忘记加锁都没机会**。从 stage05 的"靠自觉记得加锁"升级到"结构上根除忘加锁的可能"。

## 分层架构（真实 Go 项目的样子）
```
main → handler → store → model     （依赖单向，底层不认识上层）
model:   数据长什么样（Server 结构体）
store:   数据存哪、怎么安全存取（ServerStore + 锁）
handler: 处理 HTTP 请求，调 store 方法
main:    组装 & 启动（接线员）
```
> 类比 Ansible：model=defaults、store=tasks(幂等逻辑)、handler=playbook入口、main=site.yml。

## 一句话记忆
- 共享数据 + 并发请求 = **必须加锁**（stage05 的债，这关还清）
- 含锁结构体**传指针**（复制=锁失效，go vet 会抓）
- 返回数据**传副本**（append([]T(nil), s...)），别泄露内部切片
- 读多写少用 **RWMutex**；`defer Unlock()` 防死锁
- 临界区**越小越好**，IO 踢出锁外
- 锁**封装**进类型（私有字段），让人想忘都忘不了
- 最后用 **`go run -race`** 验尸，零 DATA RACE 才算过
