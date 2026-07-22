package promfmt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"metrics-collector/model"
)

// Parse 解析 Prometheus 文本格式，返回 指标名 -> 数值。
//
// 这是 stage04 的活：bufio.Scanner 逐行读、跳过注释、strconv 转数字。
// 简化点：同名指标不同标签会互相覆盖（真实 Prometheus 会按标签区分），
// 教学够用，重点是让你掌握“逐行解析文本”的手感。
func Parse(text string) (map[string]float64, error) {
	result := make(map[string]float64)
	scanner := bufio.NewScanner(strings.NewReader(text))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行 和 # 注释行（HELP / TYPE）
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 一行形如：  name{label="x"} 123   或   name 123
		// 策略：最后一个空格右边是数值，左边是“指标名(可能带标签)”
		idx := strings.LastIndex(line, " ")
		if idx < 0 {
			continue // 没有空格 = 不是合法数据行，跳过
		}
		left := strings.TrimSpace(line[:idx])
		valStr := strings.TrimSpace(line[idx+1:])

		// 指标名 = '{' 之前的部分（有标签就把标签砍掉）
		name := left
		if i := strings.Index(left, "{"); i >= 0 {
			name = left[:i]
		}

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			continue // 单行数值坏了就跳过它，别让一行脏数据毁掉整批
		}
		result[name] = val
	}
	return result, scanner.Err()
}

// WriteMetrics 把探测结果写成 Prometheus 文本格式。
//
// 注意参数是 io.Writer 而不是 http.ResponseWriter——
// 这样它既能写进 HTTP 响应，也能写进 bytes.Buffer 做单元测试，更通用。
// 这是 Go 的一个重要习惯：接口要谁的最小能力，就写谁（这里只需要“能写”）。
func WriteMetrics(w io.Writer, results []model.Result) {
	fmt.Fprintln(w, "# HELP probe_success 目标探测是否成功 (1=成功, 0=失败)")
	fmt.Fprintln(w, "# TYPE probe_success gauge")
	for _, r := range results {
		val := 0
		if r.Success {
			val = 1
		}
		// %q 自动给字符串加双引号，正好符合 label 格式 target="..."
		fmt.Fprintf(w, "probe_success{target=%q} %d\n", r.Target, val)
	}

	fmt.Fprintln(w, "# HELP probe_duration_seconds 探测耗时(秒)")
	fmt.Fprintln(w, "# TYPE probe_duration_seconds gauge")
	for _, r := range results {
		// 耗时用秒（float），Prometheus 的惯例
		fmt.Fprintf(w, "probe_duration_seconds{target=%q} %.4f\n", r.Target, r.Duration.Seconds())
	}
}
