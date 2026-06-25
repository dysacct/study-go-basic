package main

import "fmt"

func main() {
	desserts := [3]string{"Cupcake", "Eclair", "Ice cream"}
	fmt.Println(desserts)
	slice := desserts[:1] // cap will be the same as original array len
	fmt.Println("Len of slice:", len(slice), "Cap of slice:", cap(slice), "Slice:", slice)

	slice = append(slice, "Macaron")
	fmt.Println(desserts)
	fmt.Println("Len of slice:", len(slice), "Cap of slice:", cap(slice), "Slice:", slice)

	//fmt.Println(slice)
	slice = append(slice, "Cake")
	fmt.Println(desserts)
	fmt.Println("Len of slice:", len(slice), "Cap of slice:", cap(slice), "Slice:", slice)

	// bacause len is already equal to cap -> new array is allocated for the slice
	slice = append(slice, "Juice")
	fmt.Println(desserts)
	fmt.Println("Len of slice:", len(slice), "Cap of slice:", cap(slice), "Slice:", slice)

	fmt.Println("----------------")

	slice[0] = "Chocalate"
	fmt.Println(desserts)
	fmt.Println("Len of slice:", len(slice), "Cap of slice:", cap(slice), "Slice:", slice)
}
