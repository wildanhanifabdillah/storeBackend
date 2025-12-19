package services

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx         = context.Background()
	RedisClient *redis.Client
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		log.Fatal("❌ Redis connection failed:", err)
	}

	log.Println("✅ Redis connected")
}
