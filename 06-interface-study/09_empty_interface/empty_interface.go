package main

import "fmt"

type AnyValue interface{}

func LogAnyValue(v interface{}) { // type any interface{}
	fmt.Println(v)
}

// LogAnyValue == LogAnyValueWithAny
func LogAnyValueWithAny(v any) {
	fmt.Println(v)
}

func main() {
	// can assign value of any type
	var any AnyValue = "Coffee"
	fmt.Println(any)

	any = 10
	fmt.Println(any)

	any = []string{"Latte", "Espresso"}
	fmt.Println(any)

	var anotherAny interface{} = "Latte"
	anotherAny = 10.5
	anotherAny = true
	fmt.Println(anotherAny)

	// slice accepts values of an types
	var valuesOfDifferentTypes = []interface{}{
		"Latte",
		50.5,
		true,
		[3]int{1, 2, 3},
	}
	for _, v := range valuesOfDifferentTypes {
		fmt.Println(v)
	}

	// Call a function with any value
	LogAnyValue("Bogdan")
	LogAnyValue(true)
	LogAnyValue([2]string{"Latte", "Espresso"})
}
