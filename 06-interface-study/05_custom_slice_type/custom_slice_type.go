package main

import (
	"fmt"
	"strings"
)

type MenuList []string

// Stringer方法
func (ml MenuList) String() string {
	// [Coffee, Tea, Croissant]
	// option 1
	//return "[" + strings.Join(ml, ",") + "]"

	// option 2
	return fmt.Sprintf("[%s]", strings.Join(ml, ", "))

	// option 3
	//c := "["
	//for i, menuItem := range ml {
	//	c += menuItem
	//	if i < len(ml)-1 {
	//		c += ", "
	//	}
	//}
	//c += "]"
	//return c
}

func main() {
	menu := MenuList{"Coffee", "Tea", "Croissant"}
	fmt.Println("Menu:", menu)
}
