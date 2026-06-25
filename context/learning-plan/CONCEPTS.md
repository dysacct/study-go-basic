# Channel、Select、Context 核心概念图谱

> 把并发概念当作厨房动线：channel 传菜，select 盯多个窗口，context 负责喊停和截止时间。

## 本页怎么学

- 笔记文件：`context/learning-plan/CONCEPTS.md`
- 建议节奏：先跑代码，再回来看概念；并发知识只靠读很容易“懂了，但手一写就打结”。
- 观察重点：谁发送、谁接收、谁关闭、谁负责取消、谁等待退出。

## What：核心概念

这一组笔记围绕 Go 并发控制展开：`channel` 负责 goroutine 之间传递数据或信号，`select` 负责同时等待多个通信事件，`context` 负责跨函数传播取消、超时和请求级信息，`done channel` 则是轻量的完成或停止广播。

## Why：为什么要学

真实服务里，请求可能超时，用户可能断开连接，后台任务可能需要停止。如果 goroutine 只会启动不会退出，程序就会像没关水龙头的厨房，迟早把内存和资源淹掉。

## Problem：它解决什么问题

- 避免 goroutine 泄漏。
- 让多个并发任务能协同停止。
- 给慢操作设置超时边界。
- 在调用链里传递取消信号。
- 让生产者、消费者、worker pool 等模式可控退出。

## Principle：脑内模型

- `channel` 是队列加同步点：无缓冲通道要求发送者和接收者同时到场；有缓冲通道允许先存一点数据。
- `close(ch)` 是广播“没有更多数据”，不是发送一个特殊值。
- `context.Context` 是树状传播：父 context 取消，子 context 都会收到取消。
- `done := make(chan struct{})` 常用于只传信号不传数据，关闭它可以唤醒所有等待者。
- 谁创建用于发送的通道，通常谁负责关闭；接收方不要随手关别人的通道。

## Example：关联代码

### `context/learning-plan/main.go`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.After(60 * time.Second)
	c := <-t
	fmt.Println(c)
}
```

## Real World Usage：真实开发场景

- HTTP 请求超时后取消数据库查询或外部 API 调用。
- 服务收到 SIGTERM 后优雅关闭 worker 和消费者。
- 批量任务用 channel 分发，用 WaitGroup 等待完成。
- 爬虫或下载器用 semaphore 限制并发量。
- 后台任务通过 done channel 停止循环。

## Common Mistakes：高频坑点

- 发送方没人接收，goroutine 永久阻塞。
- 接收方等待一个永远不会关闭的 channel。
- 多个发送者同时 close 同一个 channel 导致 panic。
- 把 context 存进结构体长期持有，而不是按调用链传递。
- 忘记调用 `cancel()`，导致计时器和资源延迟释放。
- 用 `context.Value` 传业务参数，把 context 当杂物包。

## Practice：动手练习

- 给每个 goroutine 增加退出日志，确认它真的退出。
- 把无缓冲 channel 改成有缓冲，观察阻塞点变化。
- 给耗时任务加 `context.WithTimeout`。
- 写一个 `select` 同时等待结果和超时。
- 故意不关闭 done channel，观察程序哪里卡住。

## Interview Questions：面试问答

**问：channel 关闭后还能接收吗？**

答：可以。已缓存的数据会先被读完，之后返回元素零值和 `ok=false`。

**问：context 主要解决什么问题？**

答：跨 goroutine、跨函数调用链传播取消信号、截止时间和少量请求级值。

**问：done channel 和 context 怎么选？**

答：内部简单协调可以用 done channel；跨 API 边界、需要超时或取消原因时优先用 context。

## 一句话记忆

并发代码最重要的不是“启动多少 goroutine”，而是“每个 goroutine 怎么体面地结束”。
