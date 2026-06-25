# Go 教学重构：基本接口：Go 的鸭子类型咖啡馆

> Go 不问“你是不是官方认证咖啡机”，只问“你会不会 Brew 和 Clean”。会这两招，你就能上岗。

## 本节定位

- 笔记文件：`interface/01_basic_interface/README.md`
- 配套代码：`06-interface-study/01_basic_interface/basic_interface.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

接口是一组方法签名。某个类型只要实现了接口里的全部方法，就自动满足这个接口。

## 2. Why：为什么需要它

接口让代码依赖行为，而不是依赖具体结构体。这样替换实现、写测试、扩展功能都会轻很多。

## 3. Problem：它解决什么问题

如果函数只接受 `CapsuleMachine`，以后想支持手冲壶、意式机都得改函数。接口让函数只关心“能冲咖啡”。

## 4. Principle：底层怎么想

Go 的接口是隐式实现。接口变量内部保存动态类型和动态值。零值接口是 `nil`。

## 5. Example：本节代码

### `06-interface-study/01_basic_interface/basic_interface.go`

```go
package main

import "fmt"

type CoffeeMachine interface {
	Brew() string
	Clean() string
}

type CapsuleMachine struct {
	Brand string
}

// 满足接口类型，链接CapsuleMachine结构体
func (c CapsuleMachine) Brew() string { // 不是显式关联，而是用的相同字段名来满足接口
	return fmt.Sprintf("%s has brewed one cup of coffee", c.Brand)
}

func (c CapsuleMachine) Clean() string {
	return fmt.Sprintf("%s has cleaned", c.Brand)
}

func main() {
	var machine CoffeeMachine
	machine = CapsuleMachine{
		Brand: "Nespresso",
	}
	value := machine.Brew()
	cleaned := machine.Clean()
	fmt.Println(value)
	fmt.Println(cleaned)

}
```

## 6. Real World Usage：真实开发怎么用

- 数据库抽象
- 文件读写抽象
- 支付渠道抽象
- 测试 mock

## 7. Common Mistakes：高频坑点

- 把接口定义得过大。Go 喜欢小接口。
- 以为需要写 `implements`。Go 没有这个关键字。
- 忽略 nil 接口和装了 nil 指针的接口不是一回事。

## 8. Practice：动手练习

- 新增 `PourOverMachine` 实现同样接口。
- 写函数 `Serve(machine CoffeeMachine)` 调用 `Brew`。
- 声明一个 nil 接口并判断。

## 9. Interview Questions：面试问答

**问：Go 接口是显式实现还是隐式实现？**

答：隐式实现。只要方法集匹配即可。

**问：接口的核心价值是什么？**

答：让代码依赖行为契约，降低具体类型耦合。


## 10. 一句话记忆

接口不问你是谁，只问你会什么。
