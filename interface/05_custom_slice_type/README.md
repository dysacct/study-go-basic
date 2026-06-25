# Go 教学重构：自定义切片类型：给一串数据起个职业名

> `[]CoffeeMachine` 是“很多机器”；`CoffeeFleet` 是“咖啡机舰队”。名字一换，脑子立刻少绕两圈。

## 本节定位

- 笔记文件：`interface/05_custom_slice_type/README.md`
- 配套代码：`06-interface-study/05_custom_slice_type/custom_slice_type.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

可以用 `type CoffeeFleet []CoffeeMachine` 定义自定义切片类型，并为它添加方法。

## 2. Why：为什么需要它

当切片承载业务含义时，自定义类型能把相关操作聚合起来。

## 3. Problem：它解决什么问题

普通切片只能散落着写函数；自定义切片类型可以拥有 `TotalCapacity()`、`Filter()` 等方法。

## 4. Principle：底层怎么想

自定义类型底层仍是切片，但它是新类型，可以定义方法。切片本身包含指针、长度、容量三个部分。

## 5. Example：本节代码

### `06-interface-study/05_custom_slice_type/custom_slice_type.go`

```go
package main

import (
	"fmt"
	"strings"
)

type MenuList []string

// Stringer方法
func (ml MenuList) String() string {
	// [Coffee, Tea, Croissant]
	// option 1
	//return "[" + strings.Join(ml, ",") + "]"

	// option 2
	return fmt.Sprintf("[%s]", strings.Join(ml, ", "))

	// option 3
	//c := "["
	//for i, menuItem := range ml {
	//	c += menuItem
	//	if i < len(ml)-1 {
	//		c += ", "
	//	}
	//}
	//c += "]"
	//return c
}

func main() {
	menu := MenuList{"Coffee", "Tea", "Croissant"}
	fmt.Println("Menu:", menu)
}
```

## 6. Real World Usage：真实开发怎么用

- 订单列表统计金额
- 用户列表过滤活跃用户
- 任务列表批量执行
- 配置项集合校验

## 7. Common Mistakes：高频坑点

- 以为自定义切片和原切片完全同类型。它们需要显式转换。
- 在方法里 append 后忘记返回新切片。
- 忽略底层数组共享导致修改互相影响。

## 8. Practice：动手练习

- 为自定义切片增加 `Count()` 方法。
- 实现 `ActiveOnly()` 返回过滤后的新切片。
- 观察 append 是否改变原切片。

## 9. Interview Questions：面试问答

**问：自定义切片类型能定义方法吗？**

答：能，因为它是命名类型。

**问：切片传参会复制什么？**

答：复制切片头，底层数组仍共享。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
