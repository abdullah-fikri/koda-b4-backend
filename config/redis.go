package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Redis() {
	redisUrl := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("PASSWORD_REDIS")
	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       0,
	})
}
