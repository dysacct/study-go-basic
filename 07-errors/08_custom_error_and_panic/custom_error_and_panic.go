package main

import "fmt"

type CoffeeError string

func (c CoffeeError) Error() string {
	return string(c)
}

func main() {
	var err error
	err = CoffeeError("No coffee beans loaded!")
	if err != nil {
		fmt.Println("Error:", err)
	}
	panic(err)

}
