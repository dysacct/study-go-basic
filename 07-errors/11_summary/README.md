# Go 教学重构：错误处理总结：别让咖啡机沉默爆炸

> Go 的错误处理像店长巡场：小问题写单子，大问题拉闸，收尾动作用 defer，事故边界用 recover。

## 本节定位

- 笔记文件：`07-errors/11_summary/README.md`
- 配套代码：`07-errors/11_summary/summary.go`
- 学习目标：看懂概念、能解释原理、知道坑在哪里，并能把示例改成自己的代码。

## 1. What：它是什么

本章串起 panic、defer、recover、error、自定义错误和错误包装。

## 2. Why：为什么需要它

真实系统里失败一定会发生。成熟代码的区别不在于没有失败，而在于失败时路径清楚。

## 3. Problem：它解决什么问题

初学者容易把 panic、error、recover 混成一锅粥，最后该返回的崩溃了，该崩溃的被吞了。

## 4. Principle：底层怎么想

业务失败返回 error；资源释放用 defer；边界层可 recover；不可恢复状态才 panic；错误传递时保留上下文。

## 5. Example：本节代码

### `07-errors/11_summary/summary.go`

```go
package main

import "fmt"

func main() {
	fmt.Print("Hello World!")
}
```

## 6. Real World Usage：真实开发怎么用

- Web API 错误响应
- 数据库错误包装
- 后台任务保护
- CLI 程序退出码

## 7. Common Mistakes：高频坑点

- 忽略 error。
- 滥用 panic/recover。
- 错误信息没有上下文。
- 多层重复打印日志。

## 8. Practice：动手练习

- 把本章示例改成统一返回 error。
- 用 defer 管理资源释放。
- 写一个 recover middleware 雏形。

## 9. Interview Questions：面试问答

**问：Go 错误处理的基本模式？**

答：函数返回 error，调用者显式检查。

**问：panic、recover、defer 三者关系？**

答：panic 触发栈展开，defer 执行，defer 中 recover 可捕获 panic。


## 10. 一句话记忆

先理解它解决的问题，再记语法，Go 就不会像一盒散装螺丝。
