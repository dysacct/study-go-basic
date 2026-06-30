package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidConfig  = errors.New("invalid notifier configuration")
	ErrNetworkFailure = errors.New("network request failed")
)

type Notifier interface {
	Send(messge string) error
}

type EmailNotifier struct {
	To string
}

// 实现 Send 方法
func (e *EmailNotifier) Send(message string) error {
	if e.To == "" {
		return ErrInvalidConfig
	}

	fmt.Printf("[Email] 发送到 %s : %s\n", e.To, message)
	return nil
}

type DingTalkNotifier struct {
	Webhook string
}

func (d *DingTalkNotifier) Send(message string) error {
	if !strings.Contains(d.Webhook, "https://") {
		return ErrInvalidConfig
	}

	fmt.Printf("[DingTalk] 发送到 %s: %s\n", d.Webhook, message)
	return nil
}

type FailingNotfier struct{}

func (f *FailingNotfier) Send(message string) error {
	return ErrNetworkFailure
}

func main() {
	notifier := []Notifier{
		&EmailNotifier{
			To: "ops@example.com",
		},
		&DingTalkNotifier{
			Webhook: "https://oapi.dingtalk.com/robot/send?access_token=xxx",
		},
		&FailingNotfier{},
		&EmailNotifier{},
		&DingTalkNotifier{
			Webhook: "http://invalid",
		},
	}

	// for _, n := range notifier {
	// 	err := n.Send("服务器 CPU 使用率超过 90%")
	// 	if err != nil {
	// 		fmt.Println("发送失败:", err)
	// 	}
	// }

	SendAlert(notifier, "服务器 web-01 CPU 使用率超过 90%\"")
}

func SendAlert(notifiers []Notifier, message string) {
	for _, notifier := range notifiers {
		err := notifier.Send(message)
		if err != nil {
			if errors.Is(err, ErrInvalidConfig) {
				fmt.Printf("⚠️ 配置错误: %v\n", err)
			} else {
				fmt.Printf("❌ 发送失败: %v\n", err)
			}

			continue
		}

		fmt.Println("✅ 发送成功")
	}
}
