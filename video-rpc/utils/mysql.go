package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"video_develop-main/video-rpc/configs"
	"video_develop-main/video-rpc/global"
)

var (
	err error
)

func GlobalMysql() {
	Add := configs.AppConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Add.Mysql.User, Add.Mysql.Password, Add.Mysql.Host, Add.Mysql.Port, Add.Mysql.Database)
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.GlobalMysql.Mysql.User, config.GlobalMysql.Mysql.Password, config.GlobalMysql.Mysql.Host, config.GlobalMysql.Mysql.Port, config.GlobalMysql.Mysql.Database)
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("mysql init success")
}
