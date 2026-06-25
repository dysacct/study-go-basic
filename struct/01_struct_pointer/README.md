# Go 教学重构：结构体指针：让函数真的改到原对象

> 如果结构体是一台咖啡机，值传递就是把整台机器拍照发给函数；指针传递才是把机器钥匙递过去。照片再怎么修，厨房里的咖啡机不会变；钥匙一拧，状态就真的变了。

## 本节定位

- 笔记文件：`struct/01_struct_pointer/README.md`
- 配套代码：`05-struct_study/01_struct_pointer/struct_pointer.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

结构体指针是指向某个结构体值的地址，类型写作 `*CoffeeMachine`。函数拿到指针后，可以通过它修改原结构体字段。

## 2. Why：为什么需要它

Go 默认按值传参。结构体字段多、体积大，或者你确实想修改原数据时，传指针更合适。

## 3. Problem：它解决什么问题

如果你把结构体直接传进函数，函数得到的是副本。副本状态改成 `Out of Service`，原机器仍然一本正经地营业。

## 4. Principle：底层怎么想

`&x` 取变量地址，`*T` 表示指向 `T` 的指针。Go 访问结构体指针字段时会自动解引用，所以 `machine.Status` 等价于 `(*machine).Status`。注意：`&machine` 是“指针变量自己的地址”，不是结构体字段地址，这个点很容易看晕。

## 5. Example：本节代码

### `05-struct_study/01_struct_pointer/struct_pointer.go`

```go
package main

import "fmt"

type CoffeeMachine struct {
	Model          string
	Status         string
	OperationHours int
}

func markAsOutOfService(machine *CoffeeMachine) {
	fmt.Println("Machine status in the markAsOutOfService function:", &machine)

	machine.Status = "Out of Service"
	fmt.Println("In the function - Machine status changed to:", machine.Status)
}

func main() {
	espressoMachine := CoffeeMachine{
		Model:          "Extra calss espresso machine 234A",
		Status:         "Operational",
		OperationHours: 75,
	}
	pointerToEspressoMachine := &espressoMachine
	markAsOutOfService(pointerToEspressoMachine)
	fmt.Println("Machine status in the main function:", &pointerToEspressoMachine.Status)
	fmt.Println("Machine status in the main function:", espressoMachine.Status)
}
```

## 6. Real World Usage：真实开发怎么用

- 更新数据库实体状态
- 给配置对象打补丁
- 在方法里修改结构体内部计数器
- 避免复制较大的结构体

## 7. Common Mistakes：高频坑点

- 把 `&machine` 当成原结构体地址。`machine` 已经是指针时，`&machine` 是指针变量的地址。
- 所有结构体都无脑用指针。小结构体且不修改时，值接收者更简单。
- 返回局部变量指针时担心一定悬空。Go 会做逃逸分析，必要时放到堆上。

## 8. Practice：动手练习

- 把 `markAsOutOfService` 改成值参数，观察 `main` 里的状态是否变化。
- 新增 `AddHours(machine *CoffeeMachine, hours int)`，给运行时长累加。
- 打印 `pointerToEspressoMachine` 和 `&espressoMachine`，确认它们是否一致。

## 9. Interview Questions：面试问答

**问：值传递和指针传递的区别是什么？**

答：值传递复制数据；指针传递复制地址。后者能修改原值，也能减少大对象复制。

**问：为什么 `machine.Status` 可以直接访问字段？**

答：Go 对结构体指针字段访问提供自动解引用。


## 10. 一句话记忆

想改原对象就递钥匙，也就是指针；只给照片，函数只能修图。
