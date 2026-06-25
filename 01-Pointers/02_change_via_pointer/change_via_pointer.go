package main

import "fmt"

func main() {
	var coffeePrice = 4.50
	fmt.Println("Coffee price", coffeePrice)
	// compile Time (code you write): var coffeePrice = 4.50
	// Runtime (what machine sees) : [some memory address] = 4.50

	// compile Time (code you write): fmt.Println("Coffee price", coffeePrice)
	// Runtime (what machine sees) :  fmt.Println([some memory address], [memory address [some as step 1]])
	fmt.Println("Memory address of price 4.50", &coffeePrice)
	coffeePrice = 5.00
	fmt.Println("Memory address of price 5.00", &coffeePrice)

	// pointerToCoffeePrice := &coffeePrice // same as next line
	var pointerToCoffeePrice *float64 = &coffeePrice
	*pointerToCoffeePrice = 5.50
	fmt.Println(*pointerToCoffeePrice)
	fmt.Println("----------")
	fmt.Println("Updated coffeePrice value in memory", coffeePrice)
}
