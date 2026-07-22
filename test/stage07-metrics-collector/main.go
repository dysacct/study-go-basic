package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"metrics-collector/handler"
	"metrics-collector/prober"
	"metrics-collector/promfmt"
)

// targets 是要探测的目标清单。放这里当“配置”，真实项目里会从配置文件/环境变量读。
var targets = []string{
	"http://localhost:8080/health",  // 正常（需先启动 stage06 的服务）
	"http://localhost:8080/servers", // 正常
	"http://localhost:9999/nope",    // 故意连不上，测错误隔离
}

// main 是“接线员”：只做三件事——造零件、接线、启动。这里不写任何业务逻辑。
func main() {
	// 造零件：一个探测器（3 秒超时）
	p := prober.New(3 * time.Second)

	// ── 任务 1 & 2：启动时先命令行采集一轮，直观看到并发 + 错误隔离 ──
	fmt.Println("== 启动时先采集一轮（任务 1&2）==")
	for _, r := range p.ProbeAll(targets) {
		if r.Success {
			fmt.Printf("  ✅ %-38s status=%d  cost=%v\n", r.Target, r.StatusCode, r.Duration)
		} else {
			fmt.Printf("  ❌ %-38s FAIL  err=%s\n", r.Target, r.Err)
		}
	}

	// ── 任务 3：演示解析 Prometheus 文本 ──
	demo := `# HELP node_cpu_seconds_total CPU time
# TYPE node_cpu_seconds_total counter
node_cpu_seconds_total{mode="idle"} 12345.6
node_cpu_seconds_total{mode="user"} 678.9
node_memory_free_bytes 8589934592`
	metrics, _ := promfmt.Parse(demo)
	fmt.Println("\n== 解析 Prometheus 文本（任务 3）==")
	fmt.Printf("  解析出 %d 个指标\n", len(metrics))
	fmt.Printf("  查询 node_memory_free_bytes = %.0f\n", metrics["node_memory_free_bytes"])

	// ── 任务 4：起 exporter，接线后启动 ──
	h := handler.New(p, targets) // 把探测器和目标清单注入 handler
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", h.Metrics)

	fmt.Println("\n== Exporter 启动: http://localhost:8081/metrics（任务 4）==")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}
