package main

import "fmt"

// 定义桶结构
type Entry struct {
	key   string
	value int
}

// 定义HashMap
type HashMap struct {
	buckets [][]Entry
	size    int
}

const DefultBucketCount = 16

// 创建HashMap
func NewHashMap() *HashMap {
	return &HashMap{
		buckets: make([][]Entry, DefultBucketCount),
		size:    DefultBucketCount,
	}
}

// 简单的哈希函数
func (hm *HashMap) hash(key string) int {
	h := 0
	for _, ch := range key {
		h = (h*31 + int(ch)) % hm.size
	}
	return h
}

// 插入或更新键值对
func (hm *HashMap) Put(key string, value int) {
	index := hm.hash(key)

	// 如果这个桶还没有初始化，就先make 一个空切片
	if hm.buckets[index] == nil {
		hm.buckets[index] = make([]Entry, 0)
	}

	// 检查该桶里是否已有相同的key
	for i, entry := range hm.buckets[index] {
		if entry.key == key {
			hm.buckets[index][i].value = value
			return
		}
	}
	// 不存在，追加新的Entry
	hm.buckets[index] = append(hm.buckets[index], Entry{key: key, value: value})
}

// 查询
func (hm *HashMap) Get(key string) (int, bool) {
	index := hm.hash(key)

	if hm.buckets[index] == nil {
		return 0, false
	}

	for _, entry := range hm.buckets[index] {
		if entry.key == key {
			return entry.value, true
		}
	}
	return 0, false
}

func main() {
	hm := NewHashMap()

	hm.Put("apple", 5)
	hm.Put("banana", 8)
	hm.Put("apple", 10) // 更新 apple 的值

	if v, found := hm.Get("apple"); found {
		fmt.Println("apple ->", v)
	}

	if v, found := hm.Get("banana"); found {
		fmt.Println("banana ->", v)
	}

	if _, found := hm.Get("orange"); !found {
		fmt.Println("orange not found")
	}

	// 打印所有桶的内容，展示二维切片结构
	fmt.Println("\n所有桶的内容:")
	for i, bucket := range hm.buckets {
		if bucket != nil {
			fmt.Printf("桶 %d: %v\n", i, bucket)
		}
	}
}
