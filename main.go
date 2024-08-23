package main

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ShortenURLRequestBody struct {
	url string `json:"url"`
}

func initRouter(redisClient *redis.Client) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	router.POST("/mini", func(c *gin.Context) {
		var body ShortenURLRequestBody
		err := c.BindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		// validate URL

		// shorten URL

		// save to Redis

		// if redis is full, save to pgsql

		// return shortened URL
	})

	return router
}

func initRedis() *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return redis
}

func main() {
	redis := initRedis()
	router := initRouter(redis)

	router.Run()
}
