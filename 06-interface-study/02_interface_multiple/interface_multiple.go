package main

import "fmt"

type CoffeeMachine interface {
	Brew() string
	//DeepClean()
}

// CapsuleMachine implementation
type CapsuleMachine struct {
	Brand string
	Model string
	Price int
}

func (c CapsuleMachine) Brew() string {
	return fmt.Sprintf("%s %s has brewed a cup of capsule coffee", c.Brand, c.Model)
}

// DripMachine implementation
type DripMachine struct {
	Model string
	Price int
}

func (d DripMachine) Brew() string {
	return fmt.Sprintf("Drip coffee shot war prepared by %s", d.Model)
}

func (d DripMachine) DeepClean() {
	fmt.Println("Deep cleaning of the", d.Model)
}

func main() {
	var machineOne CoffeeMachine
	var machineTwo CoffeeMachine

	machineOne = CapsuleMachine{
		Brand: "Nespresso",
		Model: "XB23",
		Price: 100,
	}

	machineTwo = DripMachine{
		Model: "BrewPro",
		Price: 200,
	}

	fmt.Println(machineOne.Brew())
	fmt.Println(machineTwo.Brew())
	// machineTwo.DeepClean()  !!! Not possible because type of machineTwo is CoffeeMachine

	var machineThree DripMachine
	machineThree = DripMachine{
		Model: "SuperPowerDrip",
		Price: 300,
	}
	machineThree.DeepClean() // !!! Here is possible because now it has type is Dripmachine

}
