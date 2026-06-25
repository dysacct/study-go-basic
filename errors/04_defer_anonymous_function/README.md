# Go 教学重构：defer 匿名函数：临走前还能看一眼现场

> 普通 defer 像预填好的离店清单；匿名函数 defer 像临走前现场检查，能看到最新情况。

## 本节定位

- 笔记文件：`errors/04_defer_anonymous_function/README.md`
- 配套代码：`07-errors/04_defer_anonymous_function/defer_anonymous_function.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

`defer func(){ ... }()` 可以延迟执行一段匿名函数逻辑。

## 2. Why：为什么需要它

当清理动作需要读取最新变量、修改命名返回值、统一收尾时，匿名函数更灵活。

## 3. Problem：它解决什么问题

普通 defer 的参数会立即求值，无法天然拿到函数结束前的最新变量。

## 4. Principle：底层怎么想

匿名函数闭包可以捕获外部变量。捕获的是变量本身，不是当时的值；但如果把变量作为参数传给匿名函数，参数会立即求值。

## 5. Example：本节代码

### `07-errors/04_defer_anonymous_function/defer_anonymous_function.go`

```go
package main

import "fmt"

func main() {
	defer func() {
		fmt.Println("Cleaning a coffe machine...")
		fmt.Println("Suspending coffee machine...")
	}()

	defer fmt.Println("Brewing a fresh cup of espresso")
	fmt.Println("Brewing a fresh cup of cappuccino")
}
```

## 6. Real World Usage：真实开发怎么用

- 记录函数最终返回值
- recover panic
- 事务提交/回滚
- 耗时和状态统一日志

## 7. Common Mistakes：高频坑点

- 搞混闭包捕获和参数求值。
- 在 defer 里悄悄改返回值，调用者难以理解。
- 循环变量闭包捕获导致输出异常。

## 8. Practice：动手练习

- 比较 `defer fmt.Println(x)` 和 `defer func(){fmt.Println(x)}()`。
- 用 defer 匿名函数统计耗时。
- 尝试修改命名返回值。

## 9. Interview Questions：面试问答

**问：defer 匿名函数有什么优势？**

答：能执行多行逻辑，并通过闭包看到最新变量。

**问：闭包捕获的是值还是变量？**

答：通常捕获变量本身。


## 10. 一句话记忆

申请资源后立刻写 defer，未来的你会感谢现在的你。
