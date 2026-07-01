package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	log_file := "access.log"
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
	for scanner.Scan() {
		lineCount++
		fields := strings.Fields(scanner.Text())

		if len(fields) < 3 {
			continue
		}
		level := fields[2]
		levelCount[level]++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	fmt.Println("总日志行数: ", lineCount)
	fmt.Println("级别统计:")
	for level, count := range levelCount {
		fmt.Printf("  %-5s : %d\n", level, count)
	}

}
