package main

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func getRedisAddr() string {
	env := os.Getenv("ENV")

	if env == "prod" {
		return "redis:6379"
	}

	return "localhost:6379"

}

func initRedis() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: "",
		DB:       0,
	})

	return redis
}
