# Go 教学重构：结构体方法：给数据安排专属技能

> 函数像公共厨房，谁都能进；方法像咖啡机自带按钮，按钮长在机器身上，语义更自然。

## 本节定位

- 笔记文件：`struct/03_struct_methods/README.md`
- 配套代码：`05-struct_study/03_struct_methods/struct_methods.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

方法是带接收者的函数，写作 `func (m CoffeeMachine) Brew() string`。接收者让函数属于某个类型。

## 2. Why：为什么需要它

方法把数据和行为放在一起，调用者读到 `machine.Brew()` 就知道“这台机器在冲咖啡”。

## 3. Problem：它解决什么问题

全局函数多了以后，命名会拥挤，调用关系也不直观。方法能让类型自己表达能力。

## 4. Principle：底层怎么想

值接收者会复制接收者；指针接收者可以修改原对象，并避免复制。只要某个类型有对应方法，就能满足接口。

## 5. Example：本节代码

### `05-struct_study/03_struct_methods/struct_methods.go`

```go
package main

import "fmt"

type CoffeeShop struct {
	Name string
}

// 值接收器
// method with value receiver
func (shop CoffeeShop) greetShop() {
	fmt.Println("Welcome to the", shop.Name)
}

func main() {
	myShop := CoffeeShop{
		Name: "Brew & Beans",
	}
	myShop.greetShop()

}
```

## 6. Real World Usage：真实开发怎么用

- 实体对象的业务行为
- 配置对象的验证方法
- GORM 模型上的辅助方法
- 自定义类型实现 `String()`

## 7. Common Mistakes：高频坑点

- 需要修改字段却用了值接收者。
- 同一个类型的方法接收者一会儿值一会儿指针，风格混乱。
- 把所有业务逻辑都塞进方法，导致类型过胖。

## 8. Practice：动手练习

- 把一个值接收者方法改成指针接收者，观察状态变化。
- 新增 `NeedsMaintenance()` 方法，根据运行时长返回 bool。
- 让方法返回错误，而不是直接打印。

## 9. Interview Questions：面试问答

**问：什么时候用指针接收者？**

答：需要修改原对象、结构体较大、或保持方法集一致时。

**问：方法和普通函数本质一样吗？**

答：本质接近，方法只是多了接收者，调用语义更贴近类型。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
