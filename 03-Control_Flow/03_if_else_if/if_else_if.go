package main

import "fmt"

func main() {
	var points int // zero value for int is 0
	points = 75
	if points >= 100 {
		fmt.Println("Platinum member: Free coffee every day!")
	} else if points >= 50 {
		fmt.Println("Gold member: 20% discount on latte")
	} else if points >= 20 {
		fmt.Println("Silver member: Free cookie on Monday")
	} else {
		fmt.Println("Bronze member: Keep sipping to earn rewards")
	}
}
