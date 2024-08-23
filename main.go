package main

import (
	"net/url"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ShortenURLRequestBody struct {
	url string `json:"url"`
}

func generateSlug() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
		url, err := url.ParseRequestURI(body.url)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid URL",
			})
			return
		}

		// shorten URL/generate slug
		slug := generateSlug()

		// check if slug is already in use

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
