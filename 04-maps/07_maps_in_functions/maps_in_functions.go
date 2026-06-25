package main

import "fmt"

func sellCoffee(stock map[string]int, coffeeType string, quantity int) {
	currentStock, exists := stock[coffeeType]
	if exists {
		if currentStock >= quantity {
			stock[coffeeType] -= quantity
			fmt.Printf("Sold %d %s(s) cups\n", quantity, coffeeType)
			fmt.Println("sellCoffee function:")
			fmt.Println("	Stock:", stock)
			fmt.Printf("	Locatione in memory %p", &stock)

		} else {
			fmt.Printf("Not enough %s in stock. Only %d left. %d ordered\n", coffeeType, currentStock, quantity)
		}
	} else {
		fmt.Printf("%s is not available in stock.\n", coffeeType)
	}
}

func main() {
	coffeeStock := map[string]int{
		"Espresso":   10,
		"Latte":      5,
		"Cappuccino": 8,
	}

	//sellCoffee(coffeeStock, "Mocha", 1)
	sellCoffee(coffeeStock, "Espresso", 2)
	sellCoffee(coffeeStock, "Cappuccino", 4)
	sellCoffee(coffeeStock, "Latte", 6) // Not enough
	fmt.Printf("Locatione of coffeeSSSStock in memory %p", &coffeeStock)
	fmt.Println("Stock in the main function:", coffeeStock)

	fmt.Println("\nmain function:")
	fmt.Println("	Stock:", coffeeStock)
	fmt.Printf("	Locatione in memory in sellStock%p\n\n", &coffeeStock)

}
