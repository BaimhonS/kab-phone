package configs

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ConfigClients struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func SetUpConfigs() ConfigClients {
	return ConfigClients{
		DB:    ConnectDB(),
		Redis: ConnectRedis(),
	}
}
