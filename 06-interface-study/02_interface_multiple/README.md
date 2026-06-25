# Go 教学重构：多个实现：一个接口，多台机器轮流登场

> 接口像插座标准，胶囊机、滴滤壶、手冲壶都能插。插座不关心你长得帅不帅，只关心插头合不合。

## 本节定位

- 笔记文件：`06-interface-study/02_interface_multiple/README.md`
- 配套代码：`06-interface-study/02_interface_multiple/interface_multiple.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

多个不同类型可以实现同一个接口。调用方用接口变量接收它们，运行时会执行具体类型的方法。

## 2. Why：为什么需要它

这就是 Go 的多态：同一段业务逻辑，面对不同实现表现出不同结果。

## 3. Problem：它解决什么问题

没有接口时，每增加一种机器，就可能增加一套重复流程。接口把公共流程抽出来。

## 4. Principle：底层怎么想

接口调用是动态分派：接口变量记录具体类型，调用方法时执行该类型自己的实现。

## 5. Example：本节代码

### `06-interface-study/02_interface_multiple/interface_multiple.go`

```go
package main

import "fmt"

type CoffeeMachine interface {
	Brew() string
	//DeepClean()
}

// CapsuleMachine implementation
type CapsuleMachine struct {
	Brand string
	Model string
	Price int
}

func (c CapsuleMachine) Brew() string {
	return fmt.Sprintf("%s %s has brewed a cup of capsule coffee", c.Brand, c.Model)
}

// DripMachine implementation
type DripMachine struct {
	Model string
	Price int
}

func (d DripMachine) Brew() string {
	return fmt.Sprintf("Drip coffee shot war prepared by %s", d.Model)
}

func (d DripMachine) DeepClean() {
	fmt.Println("Deep cleaning of the", d.Model)
}

func main() {
	var machineOne CoffeeMachine
	var machineTwo CoffeeMachine

	machineOne = CapsuleMachine{
		Brand: "Nespresso",
		Model: "XB23",
		Price: 100,
	}

	machineTwo = DripMachine{
		Model: "BrewPro",
		Price: 200,
	}

	fmt.Println(machineOne.Brew())
	fmt.Println(machineTwo.Brew())
	// machineTwo.DeepClean()  !!! Not possible because type of machineTwo is CoffeeMachine

	var machineThree DripMachine
	machineThree = DripMachine{
		Model: "SuperPowerDrip",
		Price: 300,
	}
	machineThree.DeepClean() // !!! Here is possible because now it has type is Dripmachine

}
```

## 6. Real World Usage：真实开发怎么用

- 多种消息发送器
- 多种缓存实现
- 多数据库驱动
- 多支付平台

## 7. Common Mistakes：高频坑点

- 为了“未来可能扩展”过早抽接口。先有两个以上实现再抽通常更稳。
- 接口方法命名过泛，导致实现含义不一致。
- 返回具体类型时又暴露了实现细节。

## 8. Practice：动手练习

- 新增第三个咖啡机类型实现接口。
- 把多个机器放进 `[]CoffeeMachine` 循环调用。
- 比较接口切片和结构体切片的类型差异。

## 9. Interview Questions：面试问答

**问：什么是多态？**

答：同一接口调用，在不同具体类型上产生不同实现行为。

**问：接口切片能直接接收结构体切片吗？**

答：不能，`[]T` 和 `[]I` 是不同类型，需要逐个转换。


## 10. 一句话记忆

接口不问你是谁，只问你会什么。
