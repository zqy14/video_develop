package utils

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"public_comment/basic/config"
)

var (
	db   *gorm.DB
	Once sync.Once
	err  error
)

func InitMysql() {
	//单例模式
	Once.Do(func() {
		dsn := "root:zqy123456@tcp(14.103.243.153:3306)/zc?charset=utf8mb4&parseTime=True&loc=Local"
		global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		log.Println("mysql init success")
	})
	if err != nil {
		return
	}

	sqlBD, _ := global.DB.DB()

	//连接池参数
	sqlBD.SetMaxIdleConns(100)
	sqlBD.SetMaxOpenConns(100)
	sqlBD.SetConnMaxLifetime(time.Hour)

}

