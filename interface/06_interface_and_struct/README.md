# Go 教学重构：接口和结构体协作：数据负责状态，接口负责契约

> 结构体像演员本人，接口像角色要求。会念台词、会走位，就能出演；至于演员叫什么，导演先不管。

## 本节定位

- 笔记文件：`interface/06_interface_and_struct/README.md`
- 配套代码：`06-interface-study/06_interface_and_struct/interface_and_struct.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

结构体保存字段，方法描述行为；当方法集满足接口时，结构体就可以作为接口使用。

## 2. Why：为什么需要它

这是 Go 中组织业务代码的基本姿势：具体类型做事，接口隔离依赖。

## 3. Problem：它解决什么问题

没有接口时，高层代码依赖低层结构体，替换实现很麻烦。

## 4. Principle：底层怎么想

结构体方法集决定它是否实现接口。值接收者方法同时属于值和指针方法集；指针接收者方法只属于指针方法集。

## 5. Example：本节代码

### `06-interface-study/06_interface_and_struct/interface_and_struct.go`

```go
package main

import "fmt"

type Describable interface {
	Description() string
}

type Tea struct {
	Type string
	Size string
}

func (t Tea) Description() string {
	return fmt.Sprintf("A %s cup of %s Tea", t.Size, t.Type)
}

func main() {
	var d Describable = Tea{Type: "Green", Size: "Large"}
	fmt.Println(d.Description()) // "A Large cup of Green Tea
}
```

## 6. Real World Usage：真实开发怎么用

- service 依赖 repository 接口
- handler 依赖 usecase 接口
- 缓存层可替换 Redis/Memory
- 测试注入 fake struct

## 7. Common Mistakes：高频坑点

- 指针接收者实现接口后，却把值赋给接口变量。
- 结构体字段和接口混在一起导致依赖方向混乱。
- 接口命名只重复实现名，没有抽象行为。

## 8. Practice：动手练习

- 把方法改成指针接收者，观察赋值给接口是否报错。
- 写一个 fake 结构体满足同一接口。
- 让业务函数只依赖接口。

## 9. Interview Questions：面试问答

**问：值接收者和指针接收者对接口实现有什么影响？**

答：值接收者方法属于值和指针；指针接收者方法只属于指针。

**问：为什么接口能降低耦合？**

答：调用者只依赖行为契约，不依赖具体结构体。


## 10. 一句话记忆

接口不问你是谁，只问你会什么。
