package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"rider-device/initenal/basic/config"
)

var (
	db  *gorm.DB
	err error
)

func InitMysqlS() *gorm.DB {
	Add := config.GloBalConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Add.User, Add.Password, Add.Host, Add.Port, Add.DBName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	log.Println("mysql init success")
	return db
}
