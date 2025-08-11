package utils

import (
	"context"
	"devicemanage/devicerpc/basic/config"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	red *redis.Client
	ctx = context.Background()
)

func InitRedis() *redis.Client {
	Rds := config.GloBalConfig.Redis
	red = redis.NewClient(&redis.Options{
		Addr:     Rds.Address,
		Password: Rds.Password, // no password set
		DB:       Rds.DB,       // use default DB
	})
	err = red.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	log.Println("redis connect ok")

	return red
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	if red == nil {
		InitRedis()
	}
	return red
}

// CountDeliveries 当天配送数统计
func CountDeliveries(date string) int {
	// 模拟数据
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000) + 200
}

// CountCouriers 快递员分级统计
func CountCouriers() map[string]float64 {
	stats := map[string]int{
		"普通":  150,
		"VIP": 50,
		"钻石":  20,
	}

	total := 0
	for _, count := range stats {
		total += count
	}

	result := make(map[string]float64)
	for level, count := range stats {
		result[level] = float64(count) / float64(total) * 100
	}
	return result
}

// CountDailyOpenBoxes 使用HyperLogLog统计单日开箱数
func CountDailyOpenBoxes(date string) int64 {
	rdb := GetRedisClient()
	key := fmt.Sprintf("open_boxes:%s", date)

	// 模拟1000次开箱操作(有重复用户)
	for i := 0; i < 1000; i++ {
		userID := fmt.Sprintf("user%d", rand.Intn(500))
		rdb.PFAdd(ctx, key, userID)
	}

	count, err := rdb.PFCount(ctx, key).Result()
	if err != nil {
		log.Printf("获取开箱数失败: %v", err)
		return 0
	}
	return count
}
