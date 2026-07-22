package model

import "time"

// Result 是一次探测的结果——不管成功还是失败，都归一成一条 Result。
//
// 为什么单独放一个包、而且只有结构体没有逻辑？
// 因为它是整个程序的“公共货币”：prober 产出它、promfmt 消费它、handler 传递它。
// 谁都要用它，所以它必须站在依赖链的最底层——它不 import 任何自己人，
// 这样才不会有人反过来依赖上层，形成 import 循环。
type Result struct {
	Target     string        // 目标 URL
	Success    bool          // 是否成功（连得通 且 状态码 2xx）
	StatusCode int           // HTTP 状态码（失败时为 0）
	Duration   time.Duration // 响应耗时
	Err        string        // 错误信息（成功时为空）
}
