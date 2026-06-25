# Go 教学重构：sentinel error：给常见失败立一块路牌

> `ErrNoBeans` 比“咖啡豆没了”更像路牌：文字可以换，路牌编号别换。

## 本节定位

- 笔记文件：`07-errors/10_new_error/README.md`
- 配套代码：`07-errors/10_new_error/new_error.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

用包级变量保存可比较的固定错误，常见写法是 `var ErrX = errors.New("...")`。

## 2. Why：为什么需要它

调用者可以用 `errors.Is` 判断错误类别，而不用比较字符串。

## 3. Problem：它解决什么问题

字符串错误无法稳定判断，包装后直接 `==` 也可能失效。

## 4. Principle：底层怎么想

sentinel error 表示一类已知错误。Go 1.13 后推荐用 `errors.Is(err, ErrX)` 支持包装链。

## 5. Example：本节代码

### `07-errors/10_new_error/new_error.go`

```go
package main

import (
	"errors"
	"fmt"
)

func CheckTemperature(temp int) error {
	if temp > 120 {
		return errors.New("Critical failure: temp exceeds safe limit")
	}
	if temp > 90 {
		return fmt.Errorf("Machine overheated at %d grad C", temp)
	}
	return nil
}

func main() {
	temps := []int{75, 95, 130}

	for _, temp := range temps {
		err := CheckTemperature(temp)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("All good. Coffee machine is working fine.")
		}
	}
}
```

## 6. Real World Usage：真实开发怎么用

- `sql.ErrNoRows`
- `io.EOF`
- 业务里的 `ErrNotFound`
- 权限错误 `ErrForbidden`

## 7. Common Mistakes：高频坑点

- 导出太多 sentinel error，API 被错误变量绑死。
- 包装后还用 `==` 判断。
- 错误变量被外部修改，建议只导出 var 时谨慎。

## 8. Practice：动手练习

- 定义 `ErrInvalidCups`。
- 用 `%w` 包装它。
- 用 `errors.Is` 判断。

## 9. Interview Questions：面试问答

**问：什么是 sentinel error？**

答：预定义的固定错误值，用于表示某类错误。

**问：判断包装后的 sentinel error 用什么？**

答：`errors.Is`。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
