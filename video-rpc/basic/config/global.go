package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	GlobalMysql Config
	DB          *gorm.DB
	Reds        *redis.Client
)
