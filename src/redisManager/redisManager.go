package redisManager

import (
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	Pool *redis.Client
}

func NewRedisManager() *RedisManager {
	redis_url := os.Getenv("REDIS_URL")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis_db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Panic("REDIS_DB environment variable cannot be converted to int. Please check dotenv config...")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:           redis_url,
		Password:       redis_password,
		DB:             redis_db,
		PoolSize:       1000,
		MaxActiveConns: 0, // no limit
	})

	return &RedisManager{
		Pool: rdb,
	}
}
