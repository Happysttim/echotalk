package redis

import (
	"echotalk/internal/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisAddr,
		Password: config.Config.RedisPassword,
		DB:       config.Config.RedisDB,
	})
}
