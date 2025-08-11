package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. 当天配送数统计
	orders := []struct {
		Zone      string
		CreatedAt time.Time
	}{
		{"浦东新区", time.Now()},
		{"浦东新区", time.Now().Add(-time.Hour)},
		{"徐汇区", time.Now().AddDate(0, 0, -1)},
	}

	count := 0
	today := time.Now().Truncate(24 * time.Hour)
	for _, o := range orders {
		if o.Zone == "浦东新区" && o.CreatedAt.After(today) {
			count++
		}
	}
	fmt.Printf("当天配送数: %d\n", count)

	// 2. 快递员等级比例
	couriers := []string{"普通", "VIP", "钻石"}
	levelCount := make(map[string]int)
	for _, l := range couriers {
		levelCount[l]++
	}
	fmt.Println("快递员等级比例:")
	for level, cnt := range levelCount {
		fmt.Printf("%s: %.1f%%\n", level, float64(cnt)/float64(len(couriers))*100)
	}
}
