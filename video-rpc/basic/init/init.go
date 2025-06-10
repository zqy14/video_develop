package init

import (
	"github.com/spf13/viper"
	"log"
	"video_develop-main/video-rpc/basic/config"
)

func init() {
	InitConfig()
	//InitDB()
}

//func InitDB() {
//	config.DB = utils.GlobalMysql()
//	config.Reds = utils.ExampleClient()
//}

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("./dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		return
	}
	if err := v.Unmarshal(&config.GlobalMysql); err != nil {
		return
	}
	log.Println("viper init success")
}
