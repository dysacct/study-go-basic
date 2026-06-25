package main

import "fmt"

func main() {
	menu := map[string]float64{
		"Espresso":   2.50,
		"Latte":      3.75,
		"Cappuccino": 3.50,
		"Americano":  2.75,
	}
	fmt.Println("Original menu:", menu)
	delete(menu, "Espresso")
	fmt.Println("Update Menu is:", menu)
	delete(menu, "Coffee")
	fmt.Println("Update Menu is:", menu)
}
