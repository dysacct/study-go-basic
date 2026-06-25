# Go 教学重构：Stringer 接口：让 fmt 听懂你的类型

> `fmt.Println` 像主持人，`String()` 像嘉宾自我介绍。没准备稿，就只能念户口本字段。

## 本节定位

- 笔记文件：`06-interface-study/04_stringer_interface/README.md`
- 配套代码：`06-interface-study/04_stringer_interface/stringer_interface.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

`fmt.Stringer` 是标准库接口，要求实现 `String() string`。

## 2. Why：为什么需要它

它让自定义类型拥有统一、可读、可控的字符串展示。

## 3. Problem：它解决什么问题

日志和调试输出如果难读，排查问题就像在咖啡渣里找数据库密码。

## 4. Principle：底层怎么想

`fmt` 在格式化时会检测对象是否实现 `Stringer`，实现了就调用 `String()`。

## 5. Example：本节代码

### `06-interface-study/04_stringer_interface/stringer_interface.go`

```go
package main

import "fmt"

type Order struct {
	Customer string
	Item     string
	Quantity int
}

// 用了go默认的Stringer接口
func (o Order) String() string {
	return fmt.Sprintf("Order: %s has ordered %s (%d)", o.Customer, o.Item, o.Quantity)
}
func main() {
	order := Order{
		Customer: "Bogdan",
		Item:     "Latte",
		Quantity: 2,
	}
	fmt.Println(order.String())
}
```

## 6. Real World Usage：真实开发怎么用

- 日志摘要
- 命令行输出
- 调试复杂对象
- 错误消息中的业务对象展示

## 7. Common Mistakes：高频坑点

- `String()` 输出敏感字段。
- `String()` 中再次以字符串格式化自己导致递归。
- `String()` 太重，影响日志性能。

## 8. Practice：动手练习

- 给类型实现 `String()`。
- 使用 `fmt.Println(obj)` 验证是否自动调用。
- 给敏感字段打码。

## 9. Interview Questions：面试问答

**问：`Stringer` 属于哪个包？**

答：`fmt` 包，完整名是 `fmt.Stringer`。

**问：`String()` 常见用途？**

答：控制日志、调试和格式化输出。


## 10. 一句话记忆

接口不问你是谁，只问你会什么。
