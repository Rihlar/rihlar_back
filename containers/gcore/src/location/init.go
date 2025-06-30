package location

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	redisConn *redis.Client
)

func Init() {
	// 接続する
	conn := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	redisConn = conn
}