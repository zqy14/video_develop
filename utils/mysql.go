package utils

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
	"video_develop-main/video-rpc/configs"
	"video_develop-main/video-rpc/global"
)

func GlobalMysql() {
	var err error
	Add := configs.AppConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", Add.Mysql.User, Add.Mysql.Password, Add.Mysql.Host, Add.Mysql.Port, Add.Mysql.Database)
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.GlobalMysql.Mysql.User, config.GlobalMysql.Mysql.Password, config.GlobalMysql.Mysql.Host, config.GlobalMysql.Mysql.Port, config.GlobalMysql.Mysql.Database)
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := global.DB.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("mysql init success")
}
