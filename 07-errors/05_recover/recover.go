package main

import "fmt"

func DispenseCoffe(coffeeAmount int, cups int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Machine error:", r)
		}
	}()

	fmt.Printf("Dispensing %d grams of coffee into %d cups...", coffeeAmount, cups)
	amountPerCup := coffeeAmount / cups
	fmt.Printf("Each cup gets %d grams of coffee\n", amountPerCup)
}

func main() {
	fmt.Println("Starting coffee machine...")

	DispenseCoffe(750, 200)

	fmt.Println("Coffee machine is still running...")

	DispenseCoffe(340, 0) // error is handled using recover()

	//fmt.Println()
	fmt.Println("\nCoffee machine is still running...\n")
	DispenseCoffe(500, 150)
}
