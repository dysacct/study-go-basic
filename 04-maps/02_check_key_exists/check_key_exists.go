package main

import "fmt"

func main() {
	menu := map[string]float64{
		"Espresso": 2.50,
	}
	drink := "Cappuccino"
	fmt.Println("Cappuccino price:", menu[drink])

	price, exists := menu[drink]
	if exists {
		fmt.Printf("Price: $%.2f\n", price)
		fmt.Println("Exists:", exists)
	} else {
		fmt.Printf("%s is not on the menu\n", drink)
	}
}
