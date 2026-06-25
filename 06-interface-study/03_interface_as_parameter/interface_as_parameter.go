package main

import (
	"fmt"
)

type Barista interface {
	PrepareCoffee() string
}

type SeniorBarista struct {
	Name string
}
type JuniorBarista struct {
	Name string
}

func (s SeniorBarista) PrepareCoffee() string {
	return fmt.Sprintf("%s prepared a caramel latte", s.Name)
}

func (j JuniorBarista) PrepareCoffee() string {
	return fmt.Sprintf("%s made a hot chocolate", j.Name)
}

func ServeDrink(b Barista) {
	fmt.Println(b.PrepareCoffee())
	fmt.Println("Barista served coffee to the client")
	fmt.Println()
}

func main() {
	bogdan := SeniorBarista{Name: "Bogdan"}
	var maria Barista = JuniorBarista{Name: "Maria"}

	ServeDrink(bogdan)
	ServeDrink(maria)

	maria = SeniorBarista{Name: "Maria"}
	ServeDrink(maria)

}
