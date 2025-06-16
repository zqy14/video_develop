package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	GlobalNaCos Config
	Global      Config
	DB          *gorm.DB
	Red         *redis.Client
)
