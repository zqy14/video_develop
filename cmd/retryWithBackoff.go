package main

import (
	"fmt"
	"math"
	"time"
)

// 模拟一个可能失败的操作
func doSomething(attempt int) error {
	fmt.Printf("尝试第 %d 次操作...\n", attempt)
	// 模拟失败，前两次都失败，第三次成功
	if attempt < 3 {
		return fmt.Errorf("模拟第 %d 次操作失败", attempt)
	}
	fmt.Println("操作成功!")
	return nil
}

// 带补偿重试的操作
func retryWithBackoff(maxRetries int, operation func(int) error) error {
	var err error

	for i := 1; i <= maxRetries; i++ {
		err = operation(i)
		if err == nil {
			return nil
		}

		if i < maxRetries {
			// 计算退避时间 (2^n 秒)
			backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
			fmt.Printf("操作失败: %v. %d 秒后重试...\n", err, backoff/time.Second)

			// 等待退避时间
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("经过 %d 次重试后操作失败，最后错误: %v", maxRetries, err)
}

func main() {
	fmt.Println("开始执行带补偿重试的操作")

	err := retryWithBackoff(3, doSomething)
	if err != nil {
		fmt.Println("最终失败:", err)
	} else {
		fmt.Println("操作最终成功完成")
	}
}
