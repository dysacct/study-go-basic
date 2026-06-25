package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	t1 := time.Now()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*3))
	wg.Add(1)
	go func() {
		defer cancel()
		ip, err := GetIp(ctx, &wg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ip)
	}()
	wg.Wait()
	fmt.Println("执行完成", time.Since(t1))
}

func GetIp(ctx context.Context, wg *sync.WaitGroup) (ip string, err error) {

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("协程取消", ctx.Err())
			err = ctx.Err()
			wg.Done()
			return
		}
	}()
	time.Sleep(time.Second * 4)
	ip = "192.168.1.1"
	wg.Done()
	return
}
