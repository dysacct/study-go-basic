package main

import "fmt"

// 任务一： 结构体建模
type Server struct {
	Name   string
	CPU    int
	Memory int
}

func main() {
	s := Server{
		Name:   "server-01",
		CPU:    3,
		Memory: 4,
	}

	Upgrade(&s, 4, 6)
	fmt.Println("Server upgraded:", s)

	// 任务三：切片操作
	servers := []Server{
		{"web-01", 4, 8},
		{"web-02", 8, 16},
		{"db-01", 16, 64},
	}
	cpuTotal, memTotal := TotalResources(servers)
	fmt.Println("Total resources:", cpuTotal, memTotal)

	// 任务四：map统计
	result := CountByMemoryLevel(servers)
	fmt.Println("Count by memory level:", result)
}

// 任务二：用指针修改
func Upgrade(s *Server, addCPU, addMem int) {
	s.CPU += addCPU
	s.Memory += addMem
}

// 任务三：返回服务器的CPU和内存总和
func TotalResources(server []Server) (int, int) {
	var cpuTotal int
	var memTotal int
	for _, s := range server {
		cpuTotal += s.CPU
		memTotal += s.Memory
	}

	return cpuTotal, memTotal
}

// 任务4：map 统计
func CountByMemoryLevel(servers []Server) map[string]int {
	result := make(map[string]int)
	for _, s := range servers {
		if s.Memory <= 8 {
			result["small"]++
		} else if s.Memory <= 32 {
			result["medium"]++
		} else if s.Memory > 32 {
			result["large"]++
		}
	}
	return result
}
