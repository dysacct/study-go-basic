package main

import "fmt"

func main() {
	orderTotal := 15.0

	if orderTotal > 10 {
		fmt.Println("You get a free cookie")
	}

	orderTotal = 7.5
	if orderTotal > 10 {
		fmt.Println("You get a free cookie")
	}
}
