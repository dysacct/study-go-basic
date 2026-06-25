package main

import (
	"fmt"
	"time"
)

type Config struct {
	Addr    string
	Timeout time.Duration
}

func main() {
	cfg := Config{
		Addr:    ":8080",
		Timeout: time.Second * 5,
	}

	doRequest(cfg)
}

func doRequest(cfg Config) {
	time.Sleep(cfg.Timeout)
	fmt.Println("Request done")
}
