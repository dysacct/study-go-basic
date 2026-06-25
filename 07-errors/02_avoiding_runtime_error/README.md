# Go 教学重构：避免运行时错误：先看坑，再踩油门

> 开车前看一眼油表，不丢人；除法前看一眼分母，也不丢人。

## 本节定位

- 笔记文件：`07-errors/02_avoiding_runtime_error/README.md`
- 配套代码：`07-errors/02_avoiding_runtime_error/avoiding_runtime_error.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

通过输入校验、边界检查、nil 检查等方式，在危险操作发生前阻止 panic。

## 2. Why：为什么需要它

健壮程序不是从不遇到坏输入，而是遇到坏输入时还能体面地回答。

## 3. Problem：它解决什么问题

除零、越界、nil 解引用这类错误一旦发生，会直接中断流程。

## 4. Principle：底层怎么想

Go 鼓励显式检查。把运行时崩溃变成清晰的分支或 error，是工程代码的基本修养。

## 5. Example：本节代码

### `07-errors/02_avoiding_runtime_error/avoiding_runtime_error.go`

```go
package main

import "fmt"

func DispenseCoffe(coffeeAmount int, cups int) {
	if cups == 0 {
		fmt.Println("Error: Cannot divide coffe into 0 cups")
		return
	}
	
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

- HTTP 参数校验
- 数据库查询结果为空检查
- 切片索引前检查长度
- 配置加载后的合法性验证

## 7. Common Mistakes：高频坑点

- 只检查 happy path。
- 错误提示太模糊，定位困难。
- 校验后继续使用未校验的旧变量。

## 8. Practice：动手练习

- 给 `cups <= 0` 返回提示。
- 为切片访问写安全函数。
- 把打印错误改成返回 error。

## 9. Interview Questions：面试问答

**问：如何避免除零 panic？**

答：除法前检查分母是否为 0。

**问：Go 更推荐防御式检查还是异常捕获？**

答：Go 更推荐显式检查并返回 error。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
