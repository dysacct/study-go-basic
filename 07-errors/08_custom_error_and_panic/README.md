# Go 教学重构：自定义 error 与 panic：失败也要分级管理

> 不是每次咖啡机报警都要拉电闸。缺豆子返回 error，锅炉爆了再考虑 panic。

## 本节定位

- 笔记文件：`07-errors/08_custom_error_and_panic/README.md`
- 配套代码：`07-errors/08_custom_error_and_panic/custom_error_and_panic.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

自定义错误类型实现 `Error() string`；panic 用于严重不可恢复场景。

## 2. Why：为什么需要它

自定义 error 能携带结构化信息，调用者可以用类型断言或 `errors.As` 精准处理。

## 3. Problem：它解决什么问题

普通字符串错误表达力有限，无法可靠携带错误码、字段、业务状态。

## 4. Principle：底层怎么想

任何类型实现 `Error() string` 就满足 `error` 接口。panic 和 error 是不同层级，不应混用。

## 5. Example：本节代码

### `07-errors/08_custom_error_and_panic/custom_error_and_panic.go`

```go
package main

import "fmt"

type CoffeeError string

func (c CoffeeError) Error() string {
	return string(c)
}

func main() {
	var err error
	err = CoffeeError("No coffee beans loaded!")
	if err != nil {
		fmt.Println("Error:", err)
	}
	panic(err)

}
```

## 6. Real World Usage：真实开发怎么用

- 业务错误码
- 字段校验错误
- HTTP 状态映射
- 领域异常分类

## 7. Common Mistakes：高频坑点

- 所有失败都 panic。
- 自定义错误字段不导出，外部无法读取。
- 不用 `errors.As`，改用脆弱字符串判断。

## 8. Practice：动手练习

- 定义 `CoffeeError` 包含 Code 和 Message。
- 用 `errors.As` 识别它。
- 把业务失败从 panic 改为 error。

## 9. Interview Questions：面试问答

**问：如何定义自定义 error？**

答：定义类型并实现 `Error() string`。

**问：什么时候该 panic？**

答：不可恢复的程序错误或启动阶段致命失败。


## 10. 一句话记忆

普通失败返回 error，程序级事故才 panic。
