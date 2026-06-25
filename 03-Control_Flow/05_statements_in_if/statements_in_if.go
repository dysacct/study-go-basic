package main

import "fmt"

func getOrderWithTax(amount float64, tax float64) float64 {
	return amount + amount*tax
}
func main() {
	if points := 15; points > 10 {
		fmt.Println("You are eligible for coffee discount")
	}

	if fullAmount := getOrderWithTax(14.50, 0.1); fullAmount > 15 {
		fmt.Println("You can join coffee club...")
	} else {
		fmt.Println("You can't join coffee club...")
	}

	totalLoyaltyPoins := 150

	if totalLoyaltyPoins++; totalLoyaltyPoins > 151 {
		fmt.Println("Total loyalty points", totalLoyaltyPoins)
		fmt.Println("You can become Gold member")
	}

	if totalLoyaltyPoins += 10; totalLoyaltyPoins > 120 {
		fmt.Println("Total loyalty points", totalLoyaltyPoins)
		fmt.Println("You can become Gold member")
	}
}
