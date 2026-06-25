package main

import "fmt"

func main() {
	coffee := "Espresso"
	pointer := &coffee
	fmt.Println("Coffee name:", coffee)
	fmt.Println("Memory Location:", pointer)
	fmt.Printf("Pointer address: %p\n", pointer)

	fmt.Println("---------")
	coffeeCopy := coffee // value is copied to another location
	fmt.Println("Coffee name for coffeeCopy variable:", coffeeCopy)
	fmt.Println("Memory Location:", &coffeeCopy)
	fmt.Printf("Pointer address: %p\n", &coffeeCopy)
}
