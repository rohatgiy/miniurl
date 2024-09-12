package main

import (
	"os"

	redisbloom "github.com/RedisBloom/redisbloom-go"
	"github.com/redis/go-redis/v9"
)

func getRedisAddr() string {
	env := os.Getenv("ENV")

	if env == "prod" {
		return "redis:6379"
	}

	return "localhost:6379"

}

func initBloomFilter() *redisbloom.Client {
	bf := redisbloom.NewClient(getRedisAddr(), "bloom", nil)

	return bf
}

func initRedis() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: "",
		DB:       0,
	})

	return redis
}
