# Go 教学重构：panic：程序桌子被掀翻的那一刻

> 普通 error 像服务员说“豆子不够”；panic 像咖啡机突然喷 Steam，大厅先清场。

## 本节定位

- 笔记文件：`07-errors/01_panic_error/README.md`
- 配套代码：`07-errors/01_panic_error/panic_error.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

`panic` 表示程序遇到无法继续的严重情况。示例里 `cups` 为 0，整数除零触发运行时 panic。

## 2. Why：为什么需要它

Go 需要在越界、除零、非法内存访问等不可继续场景中立刻停止当前执行路径。

## 3. Problem：它解决什么问题

如果不提前校验输入，程序会在运行时崩溃，后面的代码不会继续执行。

## 4. Principle：底层怎么想

panic 发生后，当前 goroutine 开始栈展开，依次执行已注册的 defer；如果没有 recover，程序打印堆栈并退出。

## 5. Example：本节代码

### `07-errors/01_panic_error/panic_error.go`

```go
package main

import "fmt"

func DispenseCoffe(coffeeAmount int, cups int) {
	fmt.Printf("Dispensing %d grams of coffee into %d cups...", coffeeAmount, cups)
	amountPerCup := coffeeAmount / cups
	fmt.Printf("Each cup gets %d grams of coffee\n", amountPerCup)
}

func main() {
	fmt.Println("Starting coffee machine...")

	DispenseCoffe(750, 200)

	fmt.Println("Coffee machine is still running...")

	DispenseCoffe(340, 0) // panic: runtime error: integer divide by zero
}
```

## 6. Real World Usage：真实开发怎么用

- 程序启动阶段配置不可用可 panic
- 库内部遇到不可能状态可 panic
- 业务输入错误通常不要 panic，而是返回 error

## 7. Common Mistakes：高频坑点

- 用 panic 处理普通业务错误。
- 以为 panic 后下一行还会执行。
- 没有在危险操作前做边界检查。

## 8. Practice：动手练习

- 把 `cups` 改为 0 运行，观察堆栈。
- 在除法前增加参数校验。
- 尝试加 defer 看 panic 前是否执行。

## 9. Interview Questions：面试问答

**问：panic 和 error 的区别？**

答：error 是可预期失败；panic 是不可继续或程序缺陷。

**问：panic 时 defer 会执行吗？**

答：会，在当前 goroutine 栈展开时执行。


## 10. 一句话记忆

普通失败返回 error，程序级事故才 panic。
