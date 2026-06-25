package main

import "fmt"

func main() {
	tokens := 3

	for tokens > 0 {
		fmt.Println("Make another cup of coffee...")
		tokens--
	}

	for {
		fmt.Println("Infinite for loop")
	}
}
