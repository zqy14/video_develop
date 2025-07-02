package main

import (
	"fmt"
)

//func main() {
//	for i := 0; i < 10; i++ {
//		faker.Seed(0) // Initialize with a seed for reproducibility
//		fmt.Println(faker.Name())
//		fmt.Println(faker.Age())
//		fmt.Println(faker.IdCard())
//		fmt.Println(faker.Email())
//		fmt.Println(faker.Car())
//		fmt.Println(faker.Job())
//		fmt.Println(faker.BirthDay())
//		fmt.Println("--------------->")
//	}
//
//}
//
//package main
//
//import (
//"fmt"
//"time"
//)

//func main() {
//	// 创建一个无缓冲的channel
//	ch := make(chan string)
//
//	// 启动一个goroutine发送数据
//	go func() {
//		fmt.Println("goroutine开始发送数据...")
//		ch <- "Hello, Channel!" // 发送数据到channel
//		fmt.Println("goroutine发送数据完成")
//	}()
//
//	// 模拟主goroutine做一些工作
//	time.Sleep(1 * time.Second)
//
//	// 从channel接收数据
//	fmt.Println("主goroutine准备接收数据...")
//	msg := <-ch // 从channel接收数据
//	fmt.Println("主goroutine接收到数据:", msg)
//
//	// 等待一会儿让goroutine完成输出
//	time.Sleep(1 * time.Second)
//}
//package main

func main() {
	// 创建一个有缓冲的channel，容量为2
	ch := make(chan int, 2)

	ch <- 1
	ch <- 2

	// 从channel接收数据
	fmt.Println(<-ch) // 输出: 1
	fmt.Println(<-ch) // 输出: 2

}
