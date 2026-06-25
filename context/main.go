package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// FileInfo 存储文件的详细信息
type FileInfo struct {
	Path       string    // 文件路径
	Size       int64     // 文件大小（字节）
	ModTime    time.Time // 修改时间
	CreateTime time.Time // 创建时间（仅部分系统支持）
	LineCount  int       // 行数
}

func main() {
	// ===== 步骤1: 获取当前工作目录 =====
	// os.Getwd() 返回当前程序运行的目录路径
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前目录失败: %v\n", err)
		return
	}
	fmt.Printf("📁 开始扫描目录: %s\n", currentDir)

	// ===== 步骤2: 查找所有 .md 文件 =====
	// findMarkdownFiles 会递归遍历目录，找到所有 Markdown 文件
	mdFiles, err := findMarkdownFiles(currentDir)
	if err != nil {
		fmt.Printf("❌ 查找文件失败: %v\n", err)
		return
	}

	// ===== 步骤3: 检查是否找到文件 =====
	if len(mdFiles) == 0 {
		fmt.Println("⚠️  未找到任何 .md 文件")
		return
	}

	fmt.Printf("✅ 找到 %d 个 .md 文件，开始处理...\n\n", len(mdFiles))

	// ===== 步骤4: 创建并发控制结构 =====

	// 4.1 创建 done channel（用于通知主程序所有任务完成）
	// struct{} 是空结构体，不占用内存，用于信号传递
	done := make(chan struct{})

	// 4.2 创建结果 channel（带缓冲，避免 goroutine 阻塞）
	// 缓冲大小设为文件数量，确保所有结果都能发送
	resultChan := make(chan FileInfo, len(mdFiles))

	// 4.3 创建 WaitGroup（用于等待所有 goroutine 完成）
	// WaitGroup 是一个计数器，用于跟踪正在运行的 goroutine
	var wg sync.WaitGroup

	// ===== 步骤5: 启动并发任务 =====
	// 为每个文件启动一个 goroutine 进行处理
	for i, filePath := range mdFiles {
		wg.Add(1) // 计数器 +1，表示新增一个任务
		// 启动 goroutine（并发执行）
		go processFile(&wg, i+1, filePath, resultChan)
	}

	// ===== 步骤6: 等待所有任务完成 =====
	// 启动一个单独的 goroutine 来监控任务完成情况
	go func() {
		wg.Wait()         // 阻塞直到计数器归零（所有任务完成）
		close(resultChan) // 关闭结果 channel
		close(done)       // 关闭 done channel，通知主程序
	}()

	// ===== 步骤7: 收集处理结果 =====
	var results []FileInfo
	go func() {
		// range 会持续从 channel 接收数据，直到 channel 被关闭
		for result := range resultChan {
			results = append(results, result)
		}
	}()

	// ===== 步骤8: 等待完成信号 =====
	// 从 done channel 接收数据，阻塞直到 done 被关闭
	<-done

	// ===== 步骤9: 打印汇总信息 =====
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 处理完成！汇总信息：")
	fmt.Printf("📊 总共处理: %d 个文件\n", len(results))

	// 计算文件总大小和总行数
	var totalSize int64
	var totalLines int
	for _, result := range results {
		totalSize += result.Size
		totalLines += result.LineCount
	}
	fmt.Printf("💾 文件总大小: %.2f KB\n", float64(totalSize)/1024)
	fmt.Printf("📝 总行数: %d 行\n", totalLines)
	fmt.Println(strings.Repeat("=", 60))
}

// findMarkdownFiles 递归查找所有 .md 文件
// 参数: rootDir - 起始搜索目录
// 返回: 找到的所有 .md 文件路径列表 和 可能的错误
func findMarkdownFiles(rootDir string) ([]string, error) {
	var mdFiles []string

	// filepath.Walk 递归遍历目录树
	// 它会调用我们提供的函数，对每个文件和目录执行处理
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		// 1. 如果遍历过程中出错，直接返回错误
		if err != nil {
			return err
		}

		// 2. 跳过隐藏目录（以 . 开头的目录）
		// filepath.SkipDir 是特殊返回值，告诉 Walk 跳过这个目录
		if info.IsDir() && len(info.Name()) > 0 && info.Name()[0] == '.' {
			return filepath.SkipDir
		}

		// 3. 检查是否是 .md 文件
		// filepath.Ext() 获取文件扩展名
		// !info.IsDir() 确保不是目录
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			mdFiles = append(mdFiles, path)
		}

		// 4. 返回 nil 表示继续遍历
		return nil
	})

	return mdFiles, err
}

// processFile 处理单个文件（在独立的 goroutine 中运行）
// 参数:
//
//	wg         - WaitGroup 指针，用于通知任务完成
//	id         - Worker ID，用于日志标识
//	filePath   - 要处理的文件路径
//	resultChan - 结果 channel，用于发送处理结果
func processFile(wg *sync.WaitGroup, id int, filePath string, resultChan chan<- FileInfo) {
	// defer 确保函数退出时调用 Done()，将 WaitGroup 计数器 -1
	// 这很重要！否则主程序会一直等待
	defer wg.Done()

	// ===== 1. 获取文件相对路径（用于显示） =====
	relPath, err := filepath.Rel(".", filePath)
	if err != nil || relPath == "" {
		// 如果获取相对路径失败，使用文件名
		relPath = filepath.Base(filePath)
	}
	fmt.Printf("[Worker %d] 🔄 正在处理: %s\n", id, relPath)

	// ===== 2. 读取文件内容 =====
	// ⚠️ Go 1.16+ 使用 os.ReadFile 替代 ioutil.ReadFile
	// os.ReadFile 一次性读取整个文件到内存
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("[Worker %d] ❌ 读取文件失败 %s: %v\n", id, relPath, err)
		return
	}

	// ===== 3. 获取文件元信息 =====
	// os.Stat 返回 FileInfo 接口，包含文件的各种元数据
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("[Worker %d] ❌ 获取文件信息失败 %s: %v\n", id, relPath, err)
		return
	}

	// ===== 4. 统计文件信息 =====
	lineCount := countLines(content)
	modTime := fileInfo.ModTime() // 获取修改时间

	// ===== 5. 获取文件创建时间（系统相关） =====
	// macOS/Linux 的创建时间获取方式不同
	createTime := getFileCreationTime(fileInfo)

	// ===== 6. 打印处理结果 =====
	fmt.Printf("[Worker %d] ✅ 完成: %s\n", id, relPath)
	fmt.Printf("           大小: %d bytes | 行数: %d\n", fileInfo.Size(), lineCount)
	fmt.Printf("           修改时间: %s\n", modTime.Format("2006-01-02 15:04:05"))
	if !createTime.IsZero() {
		fmt.Printf("           创建时间: %s\n", createTime.Format("2006-01-02 15:04:05"))
	}

	// ===== 7. 发送结果到 channel =====
	// 使用 <- 操作符向 channel 发送数据
	resultChan <- FileInfo{
		Path:       filePath,
		Size:       fileInfo.Size(),
		ModTime:    modTime,
		CreateTime: createTime,
		LineCount:  lineCount,
	}
}

// countLines 统计文件行数
// 通过统计换行符 '\n' 的数量来计算行数
func countLines(content []byte) int {
	if len(content) == 0 {
		return 0
	}

	count := 0
	// 遍历每个字节，统计换行符
	for _, b := range content {
		if b == '\n' {
			count++
		}
	}

	// 如果文件最后没有换行符，也算一行
	if len(content) > 0 && content[len(content)-1] != '\n' {
		count++
	}

	return count
}

// getFileCreationTime 获取文件创建时间
// 注意: 不同操作系统的实现方式不同
// - macOS: 使用 syscall.Stat_t 的 Birthtimespec
// - Linux: 某些文件系统不支持创建时间
// - Windows: 使用不同的系统调用
func getFileCreationTime(info os.FileInfo) time.Time {
	// 尝试从系统信息中提取创建时间
	// 这部分代码依赖于操作系统
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		// macOS/BSD 系统支持 Birthtimespec
		return timespecToTime(stat.Birthtimespec)
	}

	// 如果无法获取创建时间，返回零值
	return time.Time{}
}

// timespecToTime 将 syscall.Timespec 转换为 time.Time
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(ts.Sec, ts.Nsec)
}
