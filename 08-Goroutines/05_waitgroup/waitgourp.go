package main

import (
	"fmt"
	"sync"
	"time"
)

func makeDrink(barista string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Barista %s: Starting to make coffee...\n", barista)
	time.Sleep(2. * time.Second)
	fmt.Printf("Barista %s: Done!\n", barista)
}

func main() {
	var wg sync.WaitGroup

	baristas := []string{"Bogdan", "Elena", "Maria"}
	// 主 goroutine
	fmt.Println("Coffee shop opens")

	for _, name := range baristas {
		wg.Add(1)
		go makeDrink(name, &wg)
	}
	// wg.Add(3) 也可以写在这里,因为只有三个元素,下面三行是go routine,所以可以放在for循环里
	//wg.Add(1)
	//go makeDrink("Bogdan", &wg)
	//wg.Add(1)
	//go makeDrink("Elena", &wg)
	//wg.Add(1)
	//go makeDrink("Maria", &wg)

	wg.Wait()

	fmt.Println("All drinks are ready")
	fmt.Println("Coffee shop closes")
}
