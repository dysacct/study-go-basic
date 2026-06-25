# Go 教学重构：recover：接住 panic，但别假装没事

> `recover` 像咖啡店的急停按钮：能让事故别扩散，但机器为什么冒烟还得查。

## 本节定位

- 笔记文件：`07-errors/05_recover/README.md`
- 配套代码：`07-errors/05_recover/recover.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

`recover` 只能在 defer 函数中捕获当前 goroutine 的 panic。

## 2. Why：为什么需要它

它用于边界保护，让服务、任务调度器等在局部失败后仍能继续运行。

## 3. Problem：它解决什么问题

panic 默认会让程序退出。某些边界层需要捕获它、记录日志、返回错误。

## 4. Principle：底层怎么想

panic 触发栈展开，执行 defer；defer 中调用 recover 可以停止 panic 继续传播，并得到 panic 值。

## 5. Example：本节代码

### `07-errors/05_recover/recover.go`

```go
package main

import "fmt"

func DispenseCoffe(coffeeAmount int, cups int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Machine error:", r)
		}
	}()

	fmt.Printf("Dispensing %d grams of coffee into %d cups...", coffeeAmount, cups)
	amountPerCup := coffeeAmount / cups
	fmt.Printf("Each cup gets %d grams of coffee\n", amountPerCup)
}

func main() {
	fmt.Println("Starting coffee machine...")

	DispenseCoffe(750, 200)

	fmt.Println("Coffee machine is still running...")

	DispenseCoffe(340, 0) // error is handled using recover()

	//fmt.Println()
	fmt.Println("\nCoffee machine is still running...\n")
	DispenseCoffe(500, 150)
}
```

## 6. Real World Usage：真实开发怎么用

- HTTP middleware 捕获 handler panic
- worker pool 防止单个任务炸掉进程
- 插件系统隔离第三方代码
- 测试中断言 panic

## 7. Common Mistakes：高频坑点

- 在普通函数里直接调用 recover，结果拿不到 panic。
- recover 后不记录日志，问题被吞掉。
- 以为能捕获其他 goroutine 的 panic。

## 8. Practice：动手练习

- 写一个会 panic 的函数，用 defer recover 捕获。
- 在新 goroutine 里 panic，观察外层 recover 是否有效。
- recover 后返回 error。

## 9. Interview Questions：面试问答

**问：recover 必须在哪里调用？**

答：必须在 defer 执行的函数中才有效。

**问：recover 能跨 goroutine 吗？**

答：不能，每个 goroutine 要自己保护。


## 10. 一句话记忆

普通失败返回 error，程序级事故才 panic。
