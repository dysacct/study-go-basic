# Go 教学重构：接口作为参数：函数只要能力，不查户口

> 好函数像咖啡店老板：别给我简历，能冲咖啡就来。

## 本节定位

- 笔记文件：`06-interface-study/03_interface_as_parameter/README.md`
- 配套代码：`06-interface-study/03_interface_as_parameter/interface_as_parameter.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

函数参数可以声明为接口类型，调用者传入任何实现该接口的值。

## 2. Why：为什么需要它

这样函数不会被具体结构体绑死，测试和扩展都会更自然。

## 3. Problem：它解决什么问题

函数参数写成具体类型时，复用性被限制。写成接口后，函数关注能力而非身份。

## 4. Principle：底层怎么想

参数传入接口时，具体值被装箱成接口值。方法调用通过接口的方法集完成。

## 5. Example：本节代码

### `06-interface-study/03_interface_as_parameter/interface_as_parameter.go`

```go
package main

import (
	"fmt"
)

type Barista interface {
	PrepareCoffee() string
}

type SeniorBarista struct {
	Name string
}
type JuniorBarista struct {
	Name string
}

func (s SeniorBarista) PrepareCoffee() string {
	return fmt.Sprintf("%s prepared a caramel latte", s.Name)
}

func (j JuniorBarista) PrepareCoffee() string {
	return fmt.Sprintf("%s made a hot chocolate", j.Name)
}

func ServeDrink(b Barista) {
	fmt.Println(b.PrepareCoffee())
	fmt.Println("Barista served coffee to the client")
	fmt.Println()
}

func main() {
	bogdan := SeniorBarista{Name: "Bogdan"}
	var maria Barista = JuniorBarista{Name: "Maria"}

	ServeDrink(bogdan)
	ServeDrink(maria)

	maria = SeniorBarista{Name: "Maria"}
	ServeDrink(maria)

}
```

## 6. Real World Usage：真实开发怎么用

- `io.Copy(dst Writer, src Reader)`
- 服务层依赖 repository 接口
- HTTP handler 依赖业务接口
- 单元测试传入 fake 实现

## 7. Common Mistakes：高频坑点

- 在实现方定义接口。Go 更常见是在使用方定义小接口。
- 接口参数太大，让调用者被迫实现无关方法。
- 函数内部又用类型断言到具体类型，破坏抽象。

## 8. Practice：动手练习

- 写 `MakeCoffee(m CoffeeMachine)`。
- 给测试创建一个假机器实现接口。
- 删掉函数里对具体类型的依赖。

## 9. Interview Questions：面试问答

**问：接口应该定义在使用方还是实现方？**

答：通常定义在使用方，因为使用方最清楚自己需要哪些行为。

**问：为什么 `io.Reader` 只有一个方法？**

答：小接口组合性强，适配成本低。


## 10. 一句话记忆

接口不问你是谁，只问你会什么。
