# Go 教学重构：defer：离开房间前自动关灯

> `defer` 像给函数门口贴便签：出去时记得关文件、解锁、擦咖啡渍。

## 本节定位

- 笔记文件：`07-errors/03_defer_function/README.md`
- 配套代码：`07-errors/03_defer_function/defer_function.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

`defer` 会把函数调用延迟到当前函数返回前执行。

## 2. Why：为什么需要它

资源释放逻辑和资源申请逻辑放在一起，代码更不容易漏。

## 3. Problem：它解决什么问题

没有 defer 时，多分支返回容易忘记关闭文件、释放锁、回滚事务。

## 4. Principle：底层怎么想

defer 参数会立刻求值；多个 defer 按后进先出顺序执行；panic 栈展开时 defer 也会执行。

## 5. Example：本节代码

### `07-errors/03_defer_function/defer_function.go`

```go
package main

import "fmt"

func closeShop() {
	fmt.Println("Closing the coffee shop...")
}
func main() {
	defer closeShop()
	fmt.Println("Opening the coffee shop...")
	fmt.Println("Serving a customer...")
}
```

## 6. Real World Usage：真实开发怎么用

- `defer file.Close()`
- `defer mu.Unlock()`
- `defer rows.Close()`
- 函数耗时统计

## 7. Common Mistakes：高频坑点

- 在循环里大量 defer 导致资源迟迟不释放。
- 以为 defer 参数在最后才求值。
- 忽略 defer 返回的错误，例如 Close 失败。

## 8. Practice：动手练习

- 写两个 defer，观察执行顺序。
- 在 defer 中打印变量，理解求值时机。
- 用 defer 模拟资源释放。

## 9. Interview Questions：面试问答

**问：多个 defer 执行顺序？**

答：后进先出，像栈。

**问：defer 参数什么时候求值？**

答：注册 defer 的那一刻。


## 10. 一句话记忆

申请资源后立刻写 defer，未来的你会感谢现在的你。
