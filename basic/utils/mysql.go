package utils

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"public_comment/basic/config"
)

var (
	db  *gorm.DB
	err error
)

func GlobalMysql() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Global.Mysql.User, config.Global.Mysql.Password, config.Global.Mysql.Host, config.Global.Mysql.Port, config.Global.Mysql.Database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("mysql init success")

	return db
}
