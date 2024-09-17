package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctxTimeout).Err(); err != nil {
		log.Fatalf("error connecting to redis : %v", err)
	}

	log.Println("redis connected")

	return redisClient
}
