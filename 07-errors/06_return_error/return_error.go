package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	// "coffee_orders.txt"
	file, err := os.Open("coffee_orders.txt")
	if err != nil {
		//fmt.Println("Error: could not open coffee orders file", err)
		//return
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File doesn't exist")
		} else {
			fmt.Println("General file opening error", err)
		}
		return
	}
	fmt.Println("Successfully accessed file:", file.Name())
}
