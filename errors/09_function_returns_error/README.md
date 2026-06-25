# Go 教学重构：函数返回错误：把失败交还给调用者决策

> 函数别当独裁者，遇到问题先汇报：老板，是重试、降级，还是今天不卖拿铁？

## 本节定位

- 笔记文件：`errors/09_function_returns_error/README.md`
- 配套代码：`07-errors/09_function_returns_error/function_returns_error.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

函数返回 error 后，调用者负责判断、包装、记录或继续返回。

## 2. Why：为什么需要它

失败处理通常依赖上下文，底层函数不一定知道该怎么办。

## 3. Problem：它解决什么问题

底层函数直接退出程序或打印日志，会让上层失去控制权。

## 4. Principle：底层怎么想

错误沿调用链向上传递。每一层只在能增加有用上下文时包装错误。

## 5. Example：本节代码

### `07-errors/09_function_returns_error/function_returns_error.go`

```go
package main

import (
	"fmt"
)

type OutOfStockError struct {
	Item string
}

func (e OutOfStockError) Error() string {
	return fmt.Sprintf("%s is out of stock", e.Item)
}

func ServeDrink(item string) (string, error) {
	// Here is your freshly brewed latte/espresso/cappuccino...
	quantity := stock[item]
	if quantity == 0 {
		return "", OutOfStockError{Item: item}
	}
	stock[item]--
	return fmt.Sprintf("Here is your freshly brewed %s", item), nil
}

var stock = map[string]int{
	"Espresso":   5,
	"Latte":      0,
	"Cappuccino": 10,
}

func main() {
	message, err := ServeDrink("Espresso")
	if err != nil {
		fmt.Println("Serving failed!", err)
	} else {
		fmt.Println(message)
	}
	fmt.Println()

	massage, err := ServeDrink("Latte")
	if err != nil {
		fmt.Println("Serving failed!", err)
	} else {
		fmt.Println(massage)
	}

	massage, err = ServeDrink("Tea")
	if err != nil {
		fmt.Println("Serving failed!", err)
	} else {
		fmt.Println(massage)
	}

	fmt.Println()
	fmt.Println(stock["Espresso"])

}
```

## 6. Real World Usage：真实开发怎么用

- repository 返回数据库错误给 service
- service 转换业务错误给 handler
- handler 转 HTTP 响应
- CLI main 决定退出码

## 7. Common Mistakes：高频坑点

- 每层都重复打印同一个错误。
- 包装错误时丢失原始错误。
- 忽略错误继续执行。

## 8. Practice：动手练习

- 写三层函数逐层返回 error。
- 用 `%w` 包装错误。
- 在最外层统一打印。

## 9. Interview Questions：面试问答

**问：错误应该在哪一层打印？**

答：通常在边界层统一打印或返回响应，中间层只补上下文。

**问：包装错误用什么格式动词？**

答：`%w`。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
