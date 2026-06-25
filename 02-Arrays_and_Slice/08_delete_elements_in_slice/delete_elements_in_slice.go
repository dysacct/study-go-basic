package main

import (
	"fmt"
	"slices"
)

func deleteByIndex(index int, slice []string) []string {
	//return append(slice[:index], slice[index+1:]...)
	// slices包的Delete方法
	return slices.Delete(slice, index, index+1)
}

func main() {
	coffees := []string{"Espresso", "Latte", "Cappuccino", "Mocha"}
	fmt.Println("Original menu:", coffees)
	fmt.Println("Length is :", len(coffees), "capacity of coffeeTypes:", cap(coffees))

	// go 语言里没有delete删除关键字，只能通过创建新切片然后覆盖原切片
	indexToRemove := 1
	coffees = append(coffees[:indexToRemove], coffees[indexToRemove+1:]...)
	fmt.Println("After removing element at index 1:", coffees)
	fmt.Println("Length is :", len(coffees), "capacity of coffeeTypes:", cap(coffees))

	indexToRemove = 0
	coffees = deleteByIndex(0, coffees)
	fmt.Println("After removing element at index 1:", coffees)
	fmt.Println("Length is :", len(coffees), "capacity of coffeeTypes:", cap(coffees))

}
