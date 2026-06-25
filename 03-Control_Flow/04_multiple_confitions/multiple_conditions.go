package main

import "fmt"

func main() {
	// 15:00 - 17:00
	// isMember
	// orderAmount
	hours := 16
	isMember := true
	orderAmount := 13.50

	// OPTION 1
	if hours >= 15 && hours <= 17 && isMember && orderAmount > 10 {
		fmt.Println("You get 30% off")
	} else {
		fmt.Println("No Happy Hour deals available")
	}

	// OPTION 2 (Not recommended)
	//if hours >= 15 {
	//	if hours <= 17 {
	//		if isMember {
	//			if orderAmount >= 10 {
	//				fmt.Println("You get 30% off")
	//			}
	//		}
	//	}
	//} else {
	//	fmt.Println("No Happy Hour deals available")
	//}
}
