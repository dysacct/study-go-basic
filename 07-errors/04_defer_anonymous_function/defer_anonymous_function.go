package main

import "fmt"

func main() {
	defer func() {
		fmt.Println("Cleaning a coffe machine...")
		fmt.Println("Suspending coffee machine...")
	}()

	defer fmt.Println("Brewing a fresh cup of espresso")
	fmt.Println("Brewing a fresh cup of cappuccino")
}
