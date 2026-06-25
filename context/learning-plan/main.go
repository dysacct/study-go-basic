package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.After(60 * time.Second)
	c := <-t
	fmt.Println(c)
}
