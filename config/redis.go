package config

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Redis() {
	redisUrl := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("PASSWORD_REDIS")
	dbRedis := os.Getenv("DB")
	dbFinal,_ := strconv.Atoi(dbRedis)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       dbFinal,
	})
}
