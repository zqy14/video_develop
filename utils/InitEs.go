package utils

import (
	"github.com/elastic/go-elasticsearch"
	"log"
	"rider-device/initenal/basic/config"
)

func InitEs() {
	var err error

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	config.Es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return
	}
	log.Println("ES连接成功")

}
