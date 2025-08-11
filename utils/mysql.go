package utils

import (
	"devicemanage/devicerpc/basic/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	db  *gorm.DB
	err error
)

func InitMysqlS() *gorm.DB {
	Add := config.GloBalConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Add.User, Add.Password, Add.Host, Add.Port, Add.DBName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("mysql connect ok")
	return db
}
