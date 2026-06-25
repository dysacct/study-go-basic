# Go 教学重构：接口切片：把不同类型放进同一个队伍

> 接口切片像咖啡机展会通行证：胶囊机、手冲壶、全自动机都能排一队，只要都符合入场能力。

## 本节定位

- 笔记文件：`06-interface-study/08_slice_interface/README.md`
- 配套代码：`06-interface-study/08_slice_interface/slice_interface.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

`[]SomeInterface` 可以存放不同的具体类型，只要它们都实现该接口。

## 2. Why：为什么需要它

当你需要统一调度不同实现时，接口切片很实用。

## 3. Problem：它解决什么问题

普通结构体切片只能放一种具体类型；接口切片能放多种实现。

## 4. Principle：底层怎么想

接口切片里的每个元素都是一个接口值，内部各自保存动态类型和动态值。`[]Concrete` 不能直接当作 `[]Interface`。

## 5. Example：本节代码

### `06-interface-study/08_slice_interface/slice_interface.go`

```go
package main

import "fmt"

type Greeter interface {
	Greet() string
}

type Customer struct {
	Name string
}

func (c Customer) Greet() string {
	return fmt.Sprintf("Customaer %s says: Hello! How are you?", c.Name)
}

type Staff struct {
	Role string
}

func (s Staff) Greet() string {
	return fmt.Sprintf("Staff (%s) says: Welcome to the Brew&Beans!", s.Role)
}

func main() {
	greeters := []Greeter{
		Customer{Name: "Bogdan"},
		Staff{Role: "Barista"},
		Customer{Name: "Elena"},
	}

	for _, g := range greeters {
		fmt.Println(g.Greet())
	}
	// "Customaer Bogdan says: Hello! How are you?"
	// "Staff (Barista) says: Welcome to the Brew&Beans!"
	// "Customaer Elena says: Hello! How are you?"

	fmt.Println()

	greeters = append(greeters, Staff{Role: "Cleaner"})
	for _, g := range greeters {
		fmt.Println(g.Greet())
	}
}
```

## 6. Real World Usage：真实开发怎么用

- 插件列表
- 任务处理器列表
- 多个通知渠道
- 多个中间件组件

## 7. Common Mistakes：高频坑点

- 试图把 `[]CapsuleMachine` 直接赋给 `[]CoffeeMachine`。
- 在循环里频繁类型断言，说明接口设计可能不够好。
- 忽略 nil 元素导致调用 panic。

## 8. Practice：动手练习

- 创建 `[]CoffeeMachine` 放入两种机器。
- 循环调用 `Brew()`。
- 尝试把具体切片赋给接口切片，理解编译错误。

## 9. Interview Questions：面试问答

**问：为什么 `[]T` 不能直接转成 `[]I`？**

答：两者内存布局和元素类型不同，需要逐个装箱成接口值。

**问：接口切片适合什么场景？**

答：统一管理多个不同实现。


## 10. 一句话记忆

接口不问你是谁，只问你会什么。
