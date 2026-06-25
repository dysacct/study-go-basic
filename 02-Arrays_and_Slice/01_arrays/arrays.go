package main

import "fmt"

func main() {
	var coffeeSizes [3]string
	fmt.Println(coffeeSizes)
	coffeeSizes[0] = "Small"
	fmt.Println(coffeeSizes)
	coffeeSizes[1] = "Medium"
	coffeeSizes[2] = "Large"
	fmt.Println(coffeeSizes)

	coffeeSizes[2] = "Extra Large"
	fmt.Println(coffeeSizes)

	fmt.Println("First element:", coffeeSizes[0])
	//fmt.Println(len(coffeeSizes))  3
	// coffeeSizes[5] = "Max Large" // 无效的 数组 索引 '5' (3 元素的数组超出界限)
}
