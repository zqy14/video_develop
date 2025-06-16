package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"video_develop-main/video-rpc/configs"
	"video_develop-main/video-rpc/global"
)

func ExampleClient() {
	global.Red = redis.NewClient(&redis.Options{
		Addr:     configs.AppConfig.Red.Addr,
		Password: configs.AppConfig.Red.Password, // no password set
		DB:       int(configs.AppConfig.Red.DB),  // use default DB
	})

	err := global.Red.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	log.Println("redis init success")

}
