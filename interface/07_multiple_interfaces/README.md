# Go 教学重构：多个接口：小能力自由组合

> 不要一上来造“万能咖啡机器人接口”。先拆成会冲、会洗、会计费的小能力，需要时再拼装。

## 本节定位

- 笔记文件：`interface/07_multiple_interfaces/README.md`
- 配套代码：`06-interface-study/07_multiple_interfaces/multiple_interfaces.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

一个类型可以同时实现多个接口；接口也可以嵌入其他接口形成组合接口。

## 2. Why：为什么需要它

小接口让类型按需满足能力，减少无关方法负担。

## 3. Problem：它解决什么问题

大接口会逼实现者写一堆没用的方法，最后大家都在返回 `not implemented`。

## 4. Principle：底层怎么想

接口的方法集可以通过嵌入合并。类型只要拥有组合后全部方法，就实现该组合接口。

## 5. Example：本节代码

### `06-interface-study/07_multiple_interfaces/multiple_interfaces.go`

```go
package main

import "fmt"

type PaymentMethod interface {
	Pay(amount float64) string
}

type CardInfoProvider interface {
	CardInfo() string
}

type GiftCard struct {
	Code    string
	Balance float64
}

func (g GiftCard) Pay(amount float64) string {
	if amount > g.Balance {
		return "Not enough balance"
	}
	return fmt.Sprintf("Paid $%.2f using gift card", amount)
}

func (g GiftCard) CardInfo() string {
	return fmt.Sprintf("Gift card code: %s | Balance: $%.2f", g.Code, g.Balance)
}

func main() {
	card := GiftCard{Code: "GC0001", Balance: 125.0}

	var pay PaymentMethod = card
	var info CardInfoProvider = card

	fmt.Println(info.CardInfo())
	fmt.Println(pay.Pay(35.50))
}
```

## 6. Real World Usage：真实开发怎么用

- `io.Reader`、`io.Writer`、`io.ReadWriter`
- 可启动/可停止服务组件
- 可验证/可保存模型
- 缓存读写接口拆分

## 7. Common Mistakes：高频坑点

- 接口过大。
- 接口组合层级太深，调用者不知道真正需要什么。
- 为了复用而组合，结果让实现者负担变重。

## 8. Practice：动手练习

- 拆出 `Brewer` 和 `Cleaner` 两个小接口。
- 组合成 `CoffeeMachine`。
- 写函数只接收 `Brewer`，验证更灵活。

## 9. Interview Questions：面试问答

**问：Go 标准库哪个组合接口最经典？**

答：`io.ReadWriter` 组合了 `Reader` 和 `Writer`。

**问：为什么推荐小接口？**

答：更容易实现、测试和组合。


## 10. 一句话记忆

接口不问你是谁，只问你会什么。
