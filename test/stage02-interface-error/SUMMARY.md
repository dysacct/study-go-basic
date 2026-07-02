# Stage 02 总结：接口与错误处理

## 这一关做了什么
写了一个**运维告警通知系统**：定义统一的 `Notifier` 接口，实现邮件/钉钉/失败三种通知器，用自定义错误 + `errors.Is` 区分错误类型。

## 核心知识点

### 1. 接口（interface）—— 只定义"能做什么"，不管"怎么做"
```go
type Notifier interface {
    Send(message string) error
}
```
只要一个类型实现了 `Send(string) error` 方法，它**自动就是** Notifier，不需要显式声明 `implements`。
> 类比：容器运行时接口（CRI）——只要你实现了那套方法，Docker、containerd、CRI-O 都能被 k8s 调用。

### 2. 隐式实现（Go 的特色）
```go
func (e *EmailNotifier) Send(message string) error { ... }
// EmailNotifier 没写"我实现了 Notifier"，但它就是 Notifier
```
Java 要 `implements Notifier`，Go 不用——**方法齐了就算实现**（鸭子类型）。

### 3. 多态：一个切片装不同实现
```go
notifiers := []Notifier{
    &EmailNotifier{To: "..."},
    &DingTalkNotifier{Webhook: "..."},
    &FailingNotifier{},
}
for _, n := range notifiers {
    n.Send(msg)   // 同一个调用，跑不同的实现
}
```

### 4. 哨兵错误（sentinel error）—— 包级别的错误变量
```go
var (
    ErrInvalidConfig  = errors.New("invalid notifier configuration")
    ErrNetworkFailure = errors.New("network request failed")
)
```
预定义好错误，方便调用方**用 `errors.Is` 精确判断是哪种错**。

### 5. errors.Is —— 判断错误类型
```go
if errors.Is(err, ErrInvalidConfig) {
    // 是配置错误
}
```
> 类比：Shell 里判断 `$?` 退出码——`errors.Is` 就是 Go 版的"这个错误是不是那个已知的错误"。

## 关键思考题（务必记住）

**为什么切片里要用 `&EmailNotifier{}`（指针）而不是 `EmailNotifier{}`（值）？**
- 因为 `Send` 方法是**指针接收者** `func (e *EmailNotifier)`。
- 只有 `*EmailNotifier` 才满足 Notifier 接口，`EmailNotifier`（值）不满足。
- 接口内部存的是"实现者的指针 + 类型信息"。

## 加分技巧：错误包装（wrap）
```go
return fmt.Errorf("failed to connect: %w", ErrNetworkFailure)
```
`%w` 既保留了原始错误（`errors.Is` 还能识别），又加上了上下文信息。
> 类比：报错时既说"网络失败"，又补一句"在连接 xxx 时"。

## 一句话记忆
- 接口 = 一组方法签名，**方法齐了自动实现**
- 指针接收者的类型，进接口切片要用 **`&`**
- 用 **`errors.Is`** 判断错误种类，用 **`%w`** 包装错误
