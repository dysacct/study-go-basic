package main

import "fmt"

func main() {
	var coffeeTypes = [3]string{"Espresso", "Latte", "Cappuccino"}
	fmt.Println("Types of coffee:", coffeeTypes)
	fmt.Println("Length of the array:", len(coffeeTypes))

	coffeeTypes[len(coffeeTypes)-1] = "Milk"
	fmt.Println("Type of coffee:", coffeeTypes)
}
