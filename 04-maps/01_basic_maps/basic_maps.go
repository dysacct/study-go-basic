package main

import "fmt"

func main() {
	menu := map[string]float64{
		"Espresso":   2.50,
		"Latte":      3.75,
		"Cappuccino": 3.50,
		"Americano":  2.75,
	}

	fmt.Println("Menu today is:", menu)
	fmt.Printf("Latte costs: $%.2f\n", menu["Latte"])
	fmt.Printf("Americano costs: $%.2f\n", menu["Americano"])

	// 更改Latte
	menu["Latte"] = 4.25
	fmt.Printf("New Latte costs: $%.2f\n", menu["Latte"])

	fmt.Println("Menu items quantity:", len(menu))
	// 添加新的menu
	menu["Mocha"] = 4.00
	fmt.Println("Update Menu today is:", menu)
	fmt.Println("Update Menu items quantity:", len(menu))

}
