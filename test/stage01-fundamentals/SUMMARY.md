# Stage 01 总结：基础（指针 / 切片 / map / 结构体）

## 这一关做了什么
写了一个**服务器资源统计程序**：用 struct 表示服务器，用切片存多台，用 map 做聚合统计，用指针方法修改资源。

## 核心知识点

### 1. 结构体（struct）—— 运维视角的"配置对象"
```go
type Server struct {
    Name string
    CPU  int
    Mem  int
}
```
类比：就像一个 YAML/JSON 配置块，把相关字段打包在一起。

### 2. 值接收者 vs 指针接收者（重点！）
```go
func (s Server) Show()      {}   // 值接收者：拿到的是副本，改了不影响原对象
func (s *Server) Upgrade()  {}   // 指针接收者：拿到的是本体，改动生效
```
**规则**：要修改对象本身 → 用指针接收者 `*Server`。
类比：值接收者像 `cp file && 改副本`，指针接收者像直接 `vim file`。

### 3. slice（切片）—— 动态数组
```go
servers := []Server{}
servers = append(servers, Server{...})   // 追加
for _, s := range servers { ... }        // 遍历
```
类比：Shell 里的数组，但能自动扩容。

### 4. map —— 键值统计利器
```go
result := make(map[string]int)   // ⚠️ 必须 make，否则是 nil map
result["web"]++                  // 计数
```
类比：`awk` 的关联数组 `arr[key]++`。

## 踩过的坑（血泪教训）

### 坑 1：nil map 直接写入 → panic
```go
var result map[string]int   // ❌ 只声明没初始化，是 nil
result["a"] = 1             // panic: assignment to entry in nil map
```
**修复**：`result := make(map[string]int)`
**类比**：想往目录写文件，但目录还没 `mkdir`。声明 ≠ 分配空间。

### 坑 2：`=` 覆盖 vs `+=` 累加
```go
s.CPU = addCPU    // ❌ 直接覆盖，把原值冲掉了
s.CPU += addCPU   // ✅ 在原值基础上增加
```
需求说"增加"就用 `+=`，说"设置"才用 `=`。看清题意。

## 一句话记忆
- 要改对象用**指针接收者**
- map 用前必须 **make**
- 切片用 **append** 追加
