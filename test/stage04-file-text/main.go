package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	log_file := "access.log"
	report_path := "report.txt"
	file, err := os.Open(log_file)
	if err != nil {
		fmt.Println("打开文件失败: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 统计行数
	lineCount := 0

	levelCount := make(map[string]int)

	var errorDetails []string
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 3 {
			continue
		}
		level := fields[2]
		levelCount[level]++

		if level == "ERROR" {
			timeStr := fields[1]

			message := strings.Join(fields[3:], " ")
			formattedError := fmt.Sprintf("[%s] %s", timeStr, message)
			errorDetails = append(errorDetails, formattedError)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	standardLevels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	fmt.Println("总日志行数: ", lineCount)
	fmt.Println("级别统计:")
	for level, count := range levelCount {
		fmt.Printf("  %-5s : %d\n", level, count)
	}

	fmt.Println("==== ERROR 日志 ====")
	for _, detail := range errorDetails {
		fmt.Println(detail)
	}
	fmt.Printf("共发现 %d 条 ERROR\n\n", len(errorDetails))

	reportFile, err := os.Create(report_path)
	if err != nil {
		fmt.Printf("错误：创建持久化报告文件失败: %v\n", err)
		return
	}
	defer reportFile.Close() // 延迟关闭报告文件

	// 构造格式化持久化文本流并写入文件
	currentDate := time.Now().Format("2006-01-02") // 使用 Go 诞生时间作为格式化基准
	fmt.Fprintln(reportFile, "日志分析报告")
	fmt.Fprintf(reportFile, "生成时间: %s\n", currentDate)
	fmt.Fprintln(reportFile, "----------------")
	for _, lvl := range standardLevels {
		fmt.Fprintf(reportFile, "%-6s: %d\n", lvl, levelCount[lvl])
	}
	fmt.Fprintf(reportFile, "总计  : %d 行\n", lineCount)

	fmt.Printf("报告已生成: %s\n", report_path)

}
