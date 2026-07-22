package handler

import (
	"net/http"

	"metrics-collector/model"
	"metrics-collector/promfmt"
)

// ProbeAller 是 handler 对“探测器”的最小要求：能给我一批结果就行。
//
// 为什么用接口而不直接依赖 *prober.Prober？
// 因为 handler 不关心是谁在探测、怎么探测的，它只关心“拿到 []Result”。
// 用接口把它俩解耦：将来想换个假探测器做测试，或换个实现，handler 一行都不用改。
// （呼应 stage02 的接口思想：面向能力编程，不面向具体类型。）
type ProbeAller interface {
	ProbeAll(targets []string) []model.Result
}

type Handler struct {
	prober  ProbeAller
	targets []string
}

// New 把“探测器”和“要探测的目标清单”注入进来。
// handler 自己不 new prober，也不写死 targets——都是外面传进来的（依赖注入）。
func New(prober ProbeAller, targets []string) *Handler {
	return &Handler{prober: prober, targets: targets}
}

// Metrics 是 /metrics 端点：每次被访问都实时并发采集一轮，按 Prometheus 格式吐出。
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	results := h.prober.ProbeAll(h.targets) // 干活外包给 prober

	// Prometheus 文本格式的标准 Content-Type
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	promfmt.WriteMetrics(w, results) // 格式化外包给 promfmt
}
