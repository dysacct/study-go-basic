package main

import "fmt"

type CoffeeMachine interface {
	Brew() string
	Clean() string
}

type CapsuleMachine struct {
	Brand string
}

// 满足接口类型，链接CapsuleMachine结构体
func (c CapsuleMachine) Brew() string { // 不是显式关联，而是用的相同字段名来满足接口
	return fmt.Sprintf("%s has brewed one cup of coffee", c.Brand)
}

func (c CapsuleMachine) Clean() string {
	return fmt.Sprintf("%s has cleaned", c.Brand)
}

func main() {
	var machine CoffeeMachine
	machine = CapsuleMachine{
		Brand: "Nespresso",
	}
	value := machine.Brew()
	cleaned := machine.Clean()
	fmt.Println(value)
	fmt.Println(cleaned)

}
