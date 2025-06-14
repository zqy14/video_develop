package configs

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     uint64
		Database string
	}
	Red struct {
		Addr     string
		Password string
		DB       uint64
	}
	Minio struct {
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
		UseSSL          bool
		BasePath        string
		BucketUrl       string
	}
}

var AppConfig Config

func Viper() {
	viper.SetConfigFile("configs/dev.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return
	}
	log.Println("viper init success")
}
