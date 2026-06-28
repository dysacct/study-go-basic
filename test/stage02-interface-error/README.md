# Stage 02 测试：接口 + 错误处理

考察你 `06-interface-study` 和 `07-errors` 的掌握程度。
在本目录新建 `main.go`，实现一个**运维告警通知系统**。

## 背景场景（运维向）

服务器出问题了要发告警，可以发邮件、钉钉、Webhook，甚至短信。  
你要设计一个**统一的通知接口**，让不同渠道都能用同一套调用方式。

---

## 需求

### 任务 1：定义接口
定义一个接口 `Notifier`，有一个方法：
```go
Send(message string) error
```

### 任务 2：实现三种通知器（都要实现 Notifier 接口）

#### 2.1 EmailNotifier（邮件通知器）
- 字段：`To string`（收件人邮箱）
- `Send` 行为：
  - 如果 `To` 是空字符串，返回自定义错误 `ErrInvalidConfig`（见任务 3）
  - 否则打印：`[Email] 发送到 <To>: <message>` 并返回 `nil`

#### 2.2 DingTalkNotifier（钉钉通知器）
- 字段：`WebhookURL string`
- `Send` 行为：
  - 如果 `WebhookURL` 不包含 `"https://"`，返回 `ErrInvalidConfig`
  - 否则打印：`[DingTalk] 发送到 <WebhookURL>: <message>` 并返回 `nil`

#### 2.3 FailingNotifier（模拟网络失败的通知器）
- 无字段
- `Send` 行为：**总是返回自定义错误** `ErrNetworkFailure`（见任务 3）

### 任务 3：定义两种自定义错误

用 `errors.New` 或 `fmt.Errorf` 定义两个**包级别的错误变量**（sentinel error）：

```go
var (
    ErrInvalidConfig   = errors.New("invalid notifier configuration")
    ErrNetworkFailure  = errors.New("network request failed")
)
```

### 任务 4：写一个统一的发送函数

```go
func SendAlert(notifiers []Notifier, message string)
```

**行为**：
- 遍历所有 notifier，调用它们的 `Send(message)`
- 如果某个返回错误：
  - 用 `errors.Is` 判断是不是 `ErrInvalidConfig`，如果是，打印 `⚠️  配置错误: <err>`
  - 否则打印 `❌ 发送失败: <err>`
- 如果成功（`err == nil`），打印 `✅ 发送成功`

**重点**：不要因为一个失败就停止，要把所有 notifier 都试一遍。

### 任务 5：在 main 里测试

创建三个通知器：
```go
notifiers := []Notifier{
    &EmailNotifier{To: "ops@example.com"},
    &DingTalkNotifier{WebhookURL: "https://oapi.dingtalk.com/robot/send?access_token=xxx"},
    &EmailNotifier{To: ""},  // 故意留空，触发配置错误
    &DingTalkNotifier{WebhookURL: "http://invalid"},  // 故意用 http，触发配置错误
    &FailingNotifier{},  // 模拟网络失败
}
```

调用 `SendAlert(notifiers, "服务器 web-01 CPU 使用率超过 90%")`。

---

## 通过标准

- 能编译、能运行
- 输出里能看到：
  - 2 次 `✅ 发送成功`（前两个配置正确的）
  - 2 次 `⚠️  配置错误`（空邮箱和 http 开头的钉钉）
  - 1 次 `❌ 发送失败`（FailingNotifier 的网络错误）
- 接口定义正确、所有类型都实现了接口
- 用 `errors.Is` 正确判断错误类型

## 加分项（可选）

1. **给 `ErrNetworkFailure` 加上下文**：让 `FailingNotifier` 返回时用 `fmt.Errorf("failed to connect: %w", ErrNetworkFailure)`，这样既能 `errors.Is` 判断，又能看到更多信息。

2. **思考题**：为什么 `notifiers` 切片里的元素要用 `&EmailNotifier{}`（指针）而不是 `EmailNotifier{}`（值）？提示：接口存的是啥？

---

写完对我说「**检查**」即可。
