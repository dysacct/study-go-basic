package main

import "fmt"

func main() {
	var err error
	err = fmt.Errorf("Some interesting coffee machine error")
	// err = "Some interesting coffee machine error" // string doesn't it

	if err == nil {
		fmt.Println("There is no error!")
	} else {
		fmt.Println("Error occurred!", err)
	}
}
