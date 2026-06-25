package main

import "fmt"

type CoffeType string

func (coffee CoffeType) Describe() {
	fmt.Println("This is delicios", coffee)
}
func main() {
	var myCoffee CoffeType = "Espresso"

	myCoffee.Describe()
}
