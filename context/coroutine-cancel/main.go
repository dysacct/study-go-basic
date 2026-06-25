package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wait = sync.WaitGroup{}

func main() {
	// 记录开始时间
	t1 := time.Now()

	// 返回上下文和取消函数两个值，当调用取消函数时，所有监听该上下文的写成都会被通知取消
	ctx, cancel := context.WithCancel(context.Background())

	// 计数
	wait.Add(1)

	// 获取ip
	go func() {
		ip, err := GetIp(ctx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ip)
	}()

	// 2秒后触发取消
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	wait.Wait()
	fmt.Println("执行完成", time.Since(t1))
}
func GetIp(c context.Context) (ip string, err error) {
	go func() {
		select {
		case <-c.Done():
			fmt.Println("协程取消", c.Err())
			err = c.Err()
			wait.Done()
			return
		}
	}()
	time.Sleep(time.Second * 4)
	ip = "192.168.1.1"
	return ip, nil
}
