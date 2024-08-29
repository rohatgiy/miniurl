package main

import "github.com/redis/go-redis/v9"

func initRedis() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return redis
}
