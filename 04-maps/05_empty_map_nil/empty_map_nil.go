package main

import "fmt"

func main() {
	var stock map[string]int // zero value for map is nil

	fmt.Println(stock)
	fmt.Printf("Location in memory %p\n ", &stock)

	if stock == nil {
		fmt.Printf("stock is nil\n")
	}
	//stock["Espresso"] = 10
	//stock["Latte"] = 25
	//
	//fmt.Println("Products in stock", stock)
}
