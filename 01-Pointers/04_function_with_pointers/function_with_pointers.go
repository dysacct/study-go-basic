package main

import "fmt"

func applyDiscount(price *float64, discountRate float64) {
	//fmt.Println("Memory location of the price in calculatePriceAfterDiscount fn:", &price)
	*price = *price - (*price * discountRate)
}

func main() {
	//5.00
	//10% 15% 20%
	var coffeePrice float64 = 5.00
	//fmt.Println("Memory location of coffeePrice in main fn:", &coffeePrice)
	var discount float64 = 0.10
	fmt.Println("Basic coffee price:", coffeePrice)

	applyDiscount(&coffeePrice, discount)
	fmt.Println("New Coffee price:", coffeePrice)

}
