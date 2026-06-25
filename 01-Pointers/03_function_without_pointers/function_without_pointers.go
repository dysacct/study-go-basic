package main

import "fmt"

func calculatePriceAfterDiscount(price float64, discountRate float64) float64 {
	newPrice := price - (price * discountRate)
	return newPrice
}

func main() {
	//5.00
	//10% 15% 20%
	var coffeePrice float64 = 5.00
	var discount float64 = 0.10
	fmt.Println("Basic coffee price:", coffeePrice)

	coffeePrice = calculatePriceAfterDiscount(coffeePrice, discount)
	fmt.Println("New Coffee price:", coffeePrice)

}
