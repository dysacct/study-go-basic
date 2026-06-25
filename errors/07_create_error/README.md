# Go 教学重构：创建 error：把失败写成一句人能懂的话

> 坏错误像“失败了”；好错误像“cups 不能为 0”。前者让人挠头，后者让人改代码。

## 本节定位

- 笔记文件：`errors/07_create_error/README.md`
- 配套代码：`07-errors/07_create_error/create_error.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

可以用 `errors.New` 或 `fmt.Errorf` 创建 error。

## 2. Why：为什么需要它

清晰错误信息能帮助调用者判断原因，也能帮助你三周后少骂三周前的自己。

## 3. Problem：它解决什么问题

没有 error 或错误信息太笼统，调用者无法区分不同失败原因。

## 4. Principle：底层怎么想

`errors.New` 创建固定文本错误；`fmt.Errorf` 可以格式化并携带上下文。

## 5. Example：本节代码

### `07-errors/07_create_error/create_error.go`

```go
package main

import "fmt"

func main() {
	var err error
	err = fmt.Errorf("Some interesting coffee machine error")
	// err = "Some interesting coffee machine error" // string doesn't it

	if err == nil {
		fmt.Println("There is no error!")
	} else {
		fmt.Println("Error occurred!", err)
	}
}
```

## 6. Real World Usage：真实开发怎么用

- 参数非法
- 状态不允许
- 资源不存在
- 权限不足

## 7. Common Mistakes：高频坑点

- 错误文本首字母大写或带句号，拼接时不自然。
- 只写 `invalid input`，不说明哪个输入。
- 用字符串比较错误内容。

## 8. Practice：动手练习

- 用 `errors.New` 返回固定错误。
- 用 `fmt.Errorf` 加入 cups 的值。
- 调用处打印错误。

## 9. Interview Questions：面试问答

**问：`errors.New` 和 `fmt.Errorf` 区别？**

答：前者固定文本；后者支持格式化和包装。

**问：为什么不推荐字符串比较错误？**

答：文本易变且脆弱，应使用 sentinel error 或 errors.Is。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
