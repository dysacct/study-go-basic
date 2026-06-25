package main

import "fmt"

func DispenseCoffe(coffeeAmount int, cups int) {
	if cups == 0 {
		fmt.Println("Error: Cannot divide coffe into 0 cups")
		return
	}
	
	fmt.Printf("Dispensing %d grams of coffee into %d cups...", coffeeAmount, cups)
	amountPerCup := coffeeAmount / cups
	fmt.Printf("Each cup gets %d grams of coffee\n", amountPerCup)
}

func main() {
	fmt.Println("Starting coffee machine...")

	DispenseCoffe(750, 200)

	fmt.Println("Coffee machine is still running...")

	DispenseCoffe(340, 0) // panic: runtime error: integer divide by zero
}
