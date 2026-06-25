# Go 教学重构：结构体里的函数字段：把行为塞进对象工具箱

> 普通字段像咖啡机的品牌和容量；函数字段像“按钮背后的动作”。按下 `start`，机器知道该怎么干活。

## 本节定位

- 笔记文件：`struct/02_struct_function_type/README.md`
- 配套代码：`05-struct_study/02_struct_function_type/struct_function_type.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

函数也可以成为结构体字段，例如 `Start func() string`。这样一个结构体既能保存数据，也能携带可替换的行为。

## 2. Why：为什么需要它

当某个行为需要按场景替换时，函数字段很方便，比如回调、策略、测试替身。

## 3. Problem：它解决什么问题

如果行为都写死在普通函数里，想换一套逻辑就容易到处改。函数字段让“做什么”可以在创建结构体时决定。

## 4. Principle：底层怎么想

函数在 Go 里是一等值，可以赋给变量、作为参数、作为返回值，也可以存进结构体。调用前要保证函数字段不是 `nil`。

## 5. Example：本节代码

### `05-struct_study/02_struct_function_type/struct_function_type.go`

```go
package main

import "fmt"

type CoffeeShop struct {
	Name  string
	Greet func(shop CoffeeShop)
}

func greetShop(shop CoffeeShop) {
	fmt.Println("Welcome to the", shop.Name)
}

func main() {
	myShop := CoffeeShop{
		Name:  "Brew & Beans",
		Greet: greetShop,
	}

	myShop.Greet(myShop)

}
```

## 6. Real World Usage：真实开发怎么用

- HTTP 中间件回调
- 任务调度器的执行函数
- 测试时注入假实现
- 不同支付渠道的策略函数

## 7. Common Mistakes：高频坑点

- 忘记初始化函数字段，调用时触发 panic。
- 函数字段太多，结构体变成散装接口。行为稳定时优先考虑方法或接口。
- 在闭包里捕获会变化的外部变量，导致结果不符合预期。

## 8. Practice：动手练习

- 给结构体增加 `Stop func() string` 并调用。
- 创建两台机器，分别注入不同的启动逻辑。
- 调用前用 `if machine.Start != nil` 保护一下。

## 9. Interview Questions：面试问答

**问：函数字段和方法有什么区别？**

答：方法绑定到类型；函数字段是某个值里的可替换函数。

**问：函数字段的主要风险是什么？**

答：零值是 nil，直接调用会 panic。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
