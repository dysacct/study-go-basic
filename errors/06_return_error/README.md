# Go 教学重构：返回 error：Go 的正经拒绝方式

> Go 不喜欢把失败藏进异常抽屉，它会把 error 摆在返回值上：来，看见它，处理它。

## 本节定位

- 笔记文件：`errors/06_return_error/README.md`
- 配套代码：`07-errors/06_return_error/return_error.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

函数可以返回 `error`，调用者通过 `if err != nil` 判断是否失败。

## 2. Why：为什么需要它

显式错误处理让失败路径清楚、可控、可测试。

## 3. Problem：它解决什么问题

只打印错误会让调用者无法决策；panic 又太重。返回 error 是多数业务失败的合适表达。

## 4. Principle：底层怎么想

`error` 是接口：`type error interface { Error() string }`。任何实现 `Error() string` 的类型都是 error。

## 5. Example：本节代码

### `07-errors/06_return_error/return_error.go`

```go
package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	// "coffee_orders.txt"
	file, err := os.Open("coffee_orders.txt")
	if err != nil {
		//fmt.Println("Error: could not open coffee orders file", err)
		//return
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File doesn't exist")
		} else {
			fmt.Println("General file opening error", err)
		}
		return
	}
	fmt.Println("Successfully accessed file:", file.Name())
}
```

## 6. Real World Usage：真实开发怎么用

- 参数校验失败
- 数据库查询失败
- 文件读取失败
- 远程 API 调用失败

## 7. Common Mistakes：高频坑点

- 返回 error 后调用者不检查。
- 错误信息缺少上下文。
- 把正常分支和错误分支混在一起导致可读性差。

## 8. Practice：动手练习

- 让咖啡分配函数返回 `(int, error)`。
- 调用处处理 `err != nil`。
- 给错误增加参数上下文。

## 9. Interview Questions：面试问答

**问：Go 为什么常用多返回值返回 error？**

答：让正常结果和失败原因同时显式表达。

**问：`error` 的本质是什么？**

答：带 `Error() string` 方法的接口。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
