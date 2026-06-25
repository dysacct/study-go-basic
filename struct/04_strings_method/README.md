# Go 教学重构：String 方法：让类型开口说人话

> 没有 `String()` 的结构体，打印出来像身份证扫描件；有了 `String()`，它终于会自我介绍了。

## 本节定位

- 笔记文件：`struct/04_strings_method/README.md`
- 配套代码：`05-struct_study/04_strings_method/strings_method.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

`String() string` 是 `fmt.Stringer` 接口要求的方法。实现它以后，`fmt.Println` 等函数会优先使用你的字符串表示。

## 2. Why：为什么需要它

调试、日志、错误信息都需要可读输出。`String()` 能把类型的展示逻辑集中管理。

## 3. Problem：它解决什么问题

直接打印结构体常常输出字段堆叠，不适合日志和用户提示。

## 4. Principle：底层怎么想

`fmt` 包会检查值是否实现了 `String() string`。实现后，格式化输出会调用它。别在 `String()` 里用会再次触发自身的格式化方式，否则可能递归。

## 5. Example：本节代码

### `05-struct_study/04_strings_method/strings_method.go`

```go
package main

import "fmt"

type CoffeType string

func (coffee CoffeType) Describe() {
	fmt.Println("This is delicios", coffee)
}
func main() {
	var myCoffee CoffeType = "Espresso"

	myCoffee.Describe()
}
```

## 6. Real World Usage：真实开发怎么用

- 日志里打印订单摘要
- 调试复杂结构体
- CLI 输出可读对象
- 错误上下文展示

## 7. Common Mistakes：高频坑点

- 在 `String()` 里调用 `fmt.Sprintf("%s", x)` 递归调用自己。
- `String()` 做数据库查询等重操作。它应该轻量、无副作用。
- 把敏感信息如密码、Token 放进 `String()` 输出。

## 8. Practice：动手练习

- 给咖啡机增加 `String()`，输出品牌、状态和运行小时数。
- 用 `%v`、`%+v`、`%#v` 对比打印效果。
- 确认 `String()` 不泄露敏感字段。

## 9. Interview Questions：面试问答

**问：`fmt.Stringer` 的签名是什么？**

答：`type Stringer interface { String() string }`。

**问：为什么 `String()` 不应该有副作用？**

答：因为它可能在日志、调试、格式化中被频繁隐式调用。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
