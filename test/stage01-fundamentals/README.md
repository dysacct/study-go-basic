# Stage 01 测试：基础（指针 / 切片 / map / 结构体）

抽查你 `01-Pointers` ~ `05-struct_study` 的真实掌握程度。
在本目录新建 `main.go`，实现下面的需求。要求**一个文件、可直接 `go run main.go` 跑通**。

## 背景场景（运维向）
你要写一个「服务器资源统计」小程序，处理一批服务器的信息。

## 需求

### 任务 1：结构体建模
定义一个结构体 `Server`，字段：
- `Name`   string（主机名）
- `CPU`    int（CPU 核数）
- `Memory` int（内存 GB）

### 任务 2：用指针修改（考察指针）
写一个函数 `Upgrade(s *Server, addCPU, addMem int)`，
给指定服务器**增加** CPU 核数和内存。要求：调用后原始数据被真正改变。

### 任务 3：切片操作（考察 slice）
有一批服务器：
```go
servers := []Server{
    {"web-01", 4, 8},
    {"web-02", 8, 16},
    {"db-01", 16, 64},
}
```
写一个函数 `TotalResources(servers []Server) (int, int)`，
返回所有服务器的 **CPU 总核数** 和 **内存总和**。

### 任务 4：map 统计（考察 map）
写一个函数 `CountByMemoryLevel(servers []Server) map[string]int`，
按内存大小分类统计数量，规则：
- 内存 <= 8       → "small"
- 8 < 内存 <= 32  → "medium"
- 内存 > 32       → "large"

返回类似 `map[string]int{"small": 1, "medium": 1, "large": 1}`。

### 任务 5：在 main 中串起来
1. 创建上面那批 servers
2. 对 `web-01` 调用 `Upgrade`，加 4 核、加 8G，并打印升级后的结果
3. 打印资源总和
4. 打印内存分级统计

## 通过标准
- 能编译、能运行、输出正确
- `Upgrade` 必须真正修改原数据（这是考点）
- 命名规范、代码整洁

## 加分项（可选）
- 给 `Server` 实现 `String()` 方法（`fmt.Stringer`），让打印更好看
- 思考：任务 3 的 `range` 里如果想直接改 `servers` 元素，该怎么写？为什么？

---

写完对我说「**检查**」即可。
