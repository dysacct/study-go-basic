package main

import (
	"flag"
	"fmt"
	"project/internal/todo"
)

func main() {
	list := todo.NewTodoList("todo.json")

	// 定义 flag
	addFlag := flag.String("add", "", "添加任务")
	listFlag := flag.Bool("list", false, "列出任务")
	doneFlag := flag.Int("done", 0, "标记任务为完成")
	deleteFlag := flag.Int("delete", 0, "删除任务")

	// 解析 flag
	flag.Parse()

	switch {
	case *addFlag != "":
		list.Add(*addFlag)
		fmt.Println("任务已添加:", *addFlag)
	case *listFlag:
		for i, t := range list.List() {
			fmt.Printf("%d, %s (done=%v)\n", i+1, t.Title, t.Done)
		}
	case *doneFlag != 0:
		err := list.Done(*doneFlag)
		if err != nil {
			fmt.Println("任务不存在:", err)
			return
		}
		fmt.Println("任务已标记为完成:", *doneFlag)
	case *deleteFlag != 0:
		err := list.Delete(*deleteFlag)
		if err != nil {
			fmt.Println("任务不存在:", err)
			return
		}
		fmt.Println("任务已删除:", *deleteFlag)
	default:
		fmt.Println("未知命令:")
		printUsage()
	}
}

func printUsage() {
	fmt.Println("用法: ")
	fmt.Println("  --add <任务内容>  添加任务")
	fmt.Println("  --list           列出任务")
	fmt.Println("  --done <任务ID>  标记任务为完成")
	fmt.Println("  --delete <任务ID>  删除任务")
}
