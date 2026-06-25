package main

import "fmt"

func main() {
	menu := []string{"Cake", "Pie"}

	fmt.Println("Initial menu:", menu)
	fmt.Println("Length:", len(menu), "capacity:", cap(menu))
	fmt.Printf("Memory location: %p\n", &menu)
	fmt.Printf("Memory location of \"Cake\": %p\n", &menu[0])

	fmt.Println("--------------")

	menu = append(menu, "Donut")
	fmt.Println("Menu after adding donut:", menu)
	fmt.Println("Length:", len(menu), "capacity:", cap(menu))
	fmt.Printf("Memory location: %p\n", &menu)
	fmt.Printf("Memory location of \"Cake\": %p\n", &menu[0])
	fmt.Println("--------------")

	menu = append(menu, "Ice cream")
	fmt.Println("Menu after adding Ice cream:", menu)
	fmt.Println("Length:", len(menu), "capacity:", cap(menu))
	fmt.Printf("Memory location: %p\n", &menu)
	fmt.Printf("Memory location of \"Cake\": %p\n", &menu[0])

	fmt.Println("--------------")
	menu = append(menu, "cream")
	fmt.Println("Menu after adding cream:", menu)
	fmt.Println("Length:", len(menu), "capacity:", cap(menu))
	fmt.Printf("Memory location: %p\n", &menu)
	fmt.Printf("Memory location of \"Cake\": %p\n", &menu[0])

	fmt.Println("--------------")
	cupSizes := make([]string, 0, 5)
	fmt.Println("Len of cupSizes:", len(cupSizes), "capacity of cupSizes:", cap(cupSizes))
	// cupSizes[0] = "Small" // panic: runtime error: index out of range [0] with length 0
	cupSizes = append(cupSizes, "Small", "Medium")
	fmt.Println("Len of cupSizes:", len(cupSizes), "capacity of cupSizes:", cap(cupSizes))

	cupSizes[0] = "Extra Small"
	fmt.Println(cupSizes)
	fmt.Println("Len of cupSizes:", len(cupSizes), "capacity of cupSizes:", cap(cupSizes))

}
