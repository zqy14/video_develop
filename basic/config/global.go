package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	GloBalConfig Config
	BD           *gorm.DB
	Red          *redis.Client
    Es           *elasticsearch.Client

)
