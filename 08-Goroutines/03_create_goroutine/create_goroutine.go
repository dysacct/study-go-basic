package main

import (
	"fmt"
	"time"
)

func barista() {
	fmt.Println("Barista: Starting to make coffee...")
	time.Sleep(2 * time.Second)
	fmt.Println("Barista: Done!") // This will not be shown because main goroutine ends earlier than barista goroutine
}

func main() {
	// 主 goroutine
	fmt.Println("Coffee shop opens")
	// create a goroutine
	// add goKeyWord
	go barista()
	time.Sleep(3 * time.Second)
	fmt.Println("Coffee shop closes")
}
