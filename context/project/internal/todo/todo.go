package todo

import (
	"encoding/json"
	"fmt"
	"os"
)

// 单个任务
type Todo struct {
	Title string
	Done  bool
}

// 任务列表
type TodoList struct {
	items []Todo
	file  string
}

// 创建一个新的TodoList
func NewTodoList(file string) *TodoList {
	t := &TodoList{
		items: make([]Todo, 0),
		file:  file,
	}
	t.Load()
	return t
}

// 添加任务
func (t *TodoList) Add(title string) {
	t.items = append(t.items, Todo{
		Title: title,
		Done:  false,
	})
	t.save()
}

// 标记任务为完成
func (t *TodoList) Done(id int) error {
	index := id - 1
	if index < 0 || index >= len(t.items) {
		return fmt.Errorf("任务不存在")
	}
	t.items[index].Done = true
	t.save()
	return nil
}

// 删除任务
func (t *TodoList) Delete(id int) error {
	index := id - 1
	if index < 0 || index >= len(t.items) {
		return fmt.Errorf("任务不存在")
	}
	t.items = append(t.items[:index], t.items[index+1:]...)
	t.save()
	return nil
}

// 获取所有任务
func (t *TodoList) List() []Todo {
	return t.items
}

// 从文件加载
func (t *TodoList) Load() {
	data, err := os.ReadFile(t.file)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &t.items)
}

func (t *TodoList) save() {
	data, _ := json.MarshalIndent(t.items, "", "  ")
	_ = os.WriteFile(t.file, data, 0644)
}
