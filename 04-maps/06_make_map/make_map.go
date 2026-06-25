package main

import "fmt"

func main() {
	stock := make(map[string]int)
	fmt.Println(stock)
	stock["Espresso"] = 10
	stock["Latte"] = 25

	fmt.Println("Products in stock", stock)
}
