package main

import "fmt"

func main() {
	fmt.Println("Welcome to the Brew&Beans Terminal!")

	for {
		var order string
		fmt.Print("Enter your coffee name (or type 'exit' to quit)：")
		fmt.Scanln(&order)
		if order == "" {
			fmt.Println("Please enter a valid order...")
			continue
		}
		if order == "exit" || order == "quit" {
			fmt.Println("Thank you for visiting Br ew&Beans!")
			break
		}
		fmt.Println("Preparing your order...")
	}
	fmt.Println("Finishing program...")
}
