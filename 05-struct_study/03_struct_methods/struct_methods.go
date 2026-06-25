package main

import "fmt"

type CoffeeShop struct {
	Name string
}

// 值接收器
// method with value receiver
func (shop CoffeeShop) greetShop() {
	fmt.Println("Welcome to the", shop.Name)
}

func main() {
	myShop := CoffeeShop{
		Name: "Brew & Beans",
	}
	myShop.greetShop()

}
