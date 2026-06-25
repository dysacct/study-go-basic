package main

import "fmt"

func main() {
	coffeeCups := 10
	for i := 0; i < coffeeCups; i++ {
		fmt.Printf("Preparing coffee cup %d\n", i+1)
	}

	fmt.Println()

	for i := coffeeCups; i >= 1; i-- {
		fmt.Printf("Preparing coffee cup %d\n", i)
	}

}
