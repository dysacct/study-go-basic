# Go 教学重构：并发入门：Go 不一定跑得更快，但会安排活

> 并发不是一群人同时冲进厨房抢锅，而是店长把洗杯、磨豆、收银这些事安排得不互相干等。

## 本节定位

- 笔记文件：`Goroutines/01_concurrency/README.md`
- 配套代码：`08-Goroutines/01_concurrency/concurrency.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

并发是把任务拆成可独立推进的多个执行流。示例里的 `runtime.GOMAXPROCS(0)` 会打印当前 Go 程序可同时执行 Go 代码的逻辑处理器数量。

## 2. Why：为什么需要它

网络请求、文件 IO、后台任务经常会等待。并发能让等待时间被其他工作填起来。

## 3. Problem：它解决什么问题

串行程序遇到慢任务就整条队伍卡住。并发让多个任务交替推进。

## 4. Principle：底层怎么想

goroutine 是 Go 的轻量级并发执行单元，由 Go runtime 调度到 OS 线程上运行。`GOMAXPROCS` 控制同时执行 Go 代码的 P 数量。

## 5. Example：本节代码

### `08-Goroutines/01_concurrency/concurrency.go`

```go
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOMAXPROCS(0))
}
```

## 6. Real World Usage：真实开发怎么用

- Web 服务器同时处理多个请求
- 批量调用外部 API
- 后台消费者处理消息
- 定时任务和主服务并行运行

## 7. Common Mistakes：高频坑点

- 把并发等同于并行。并发是结构，并行是同一时刻真的同时跑。
- 启动 goroutine 后 main 直接退出。
- 共享变量不加同步导致数据竞争。

## 8. Practice：动手练习

- 运行示例，观察 `GOMAXPROCS` 输出。
- 写两个 goroutine 分别打印任务名。
- 用 `sync.WaitGroup` 等待 goroutine 结束。

## 9. Interview Questions：面试问答

**问：并发和并行区别？**

答：并发是多个任务交替推进的设计；并行是多个任务同一时刻同时执行。

**问：goroutine 由谁调度？**

答：由 Go runtime 调度到 OS 线程。


## 10. 一句话记忆

并发是安排任务不干等，并行是任务真的同时跑。
