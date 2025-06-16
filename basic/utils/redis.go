package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"public_comment/basic/config"
)

var red *redis.Client

func ExampleClient() *redis.Client {
	red = redis.NewClient(&redis.Options{
		Addr:     config.Global.Redis.Addr,
		Password: config.Global.Redis.Password, // no password set
		DB:       int(config.Global.Redis.DB),  // use default DB
	})

	err = red.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	log.Println("redis init success")
	return red
}
