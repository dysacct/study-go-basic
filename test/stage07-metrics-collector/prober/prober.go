package prober

import (
	"io"
	"net/http"
	"sync"
	"time"

	"metrics-collector/model"
)

// Prober 是“探测器”。为什么做成结构体、把 *http.Client 攥在手里？
// 因为一个 http.Client 内部有连接池，全程序共用一个最高效
// （复用 TCP 连接，别每次探测都新建）。攥在字段里，probe 时随手就能用。
type Prober struct {
	client *http.Client
}

// New 构造 Prober。超时通过参数传进来，而不是写死——
// 谁用谁决定超时多少，这叫“配置从外面注入”。
func New(timeout time.Duration) *Prober {
	return &Prober{
		client: &http.Client{Timeout: timeout},
	}
}

// Probe 探测单个目标。这一关三条铁律全在这里：
//  1. 超时靠 client.Timeout（绝不用裸 http.Get）
//  2. resp.Body 必须 defer Close（否则连接泄漏）
//  3. err==nil ≠ 成功，还要看 StatusCode 是不是 2xx
func (p *Prober) Probe(url string) model.Result {
	start := time.Now()

	resp, err := p.client.Get(url)
	if err != nil {
		// 网络层就挂了：连不上 / 超时 / DNS 失败。
		// 关键：不 panic、不 return error 往外抛，而是“把失败也记成一条结果”。
		// 这就是错误隔离——一个目标挂掉，只影响它自己这条 Result。
		return model.Result{
			Target:   url,
			Success:  false,
			Duration: time.Since(start),
			Err:      err.Error(),
		}
	}
	defer resp.Body.Close() // 确认 err==nil（resp 非空）之后再 defer，顺序不能反

	// 即使不关心响应体内容，也要把它读干净再关，
	// 否则底层连接无法被复用（HTTP keep-alive 的要求）。
	_, _ = io.Copy(io.Discard, resp.Body)

	return model.Result{
		Target:     url,
		Success:    resp.StatusCode >= 200 && resp.StatusCode < 300,
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
}

// ProbeAll 并发探测所有目标——这就是 stage03（goroutine+WaitGroup）
// 和 stage06（Mutex 保护共享切片）的合体。
func (p *Prober) ProbeAll(targets []string) []model.Result {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make([]model.Result, 0, len(targets))
	)

	for _, url := range targets {
		wg.Add(1)
		go func(u string) { // ⭐ u 当参数传进来，每个 goroutine 拿到自己的副本
			defer wg.Done()

			r := p.Probe(u) // 慢活（网络 IO）在锁外做，别占着锁去联网

			mu.Lock()
			results = append(results, r) // 只有这一行碰共享切片，把它锁最小
			mu.Unlock()
		}(url)
	}

	wg.Wait() // 等所有目标都探测完再返回
	return results
}
