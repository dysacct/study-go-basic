package main

import "fmt"

func main() {
	menu := [3]string{"Tea", "Coffee", "Juice"}
	slice := menu[:2]
	fmt.Println("Before slice change:")
	fmt.Println("menu:", menu)
	fmt.Println("slice:", slice)

	fmt.Println("Length of menu array:", len(menu))
	fmt.Println("Length of slice:", len(slice))
	//fmt.Printf("%p: menu, %p: slice\n", &menu, &slice)

	slice[0] = "Matcha"
	fmt.Println("After slice change:")
	fmt.Println("menu:", menu) // changed
	fmt.Println("slice:", slice)

}
