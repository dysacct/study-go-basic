package main

import "fmt"

func main() {
	day := "Saturday"

	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday":
		fmt.Println("Today is Workday")
	case "Sunday", "Saturday":
		fmt.Println("This is Weekend")
	default:
		fmt.Println("Invalid day")
	}
}
