package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/redis/go-redis/v9"
)

func shortenURL(postgresClient *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body ShortenURLRequestBody
		err := c.BindJSON(&body)

		if err != nil {
			c.JSON(400, gin.H{
				"error":  "Invalid request body",
				"reason": err.Error(),
			})
			return
		}

		// validate URL
		url, err := url.ParseRequestURI(body.Url)
		if err != nil {
			c.JSON(400, gin.H{
				"error":  "Invalid URL",
				"reason": err.Error(),
			})
			return
		}

		// shorten URL/generate slug
		slug := generateSlug()

		// check if slug is already in use
		// use bloom filter?
		// add retry logic

		slugExists, err := checkIfSlugExists(postgresClient, slug)
		if slugExists {
			c.JSON(500, gin.H{
				"error":  "Slug already in use",
				"reason": "Please try again",
			})
			return
		}

		if err != pg.ErrNoRows {
			c.JSON(500, gin.H{
				"error":  "Internal server error",
				"reason": err.Error(),
			})
			return
		}

		// save to pgsql
		_, err = saveSlug(postgresClient, &Slug{
			Url:  url.String(),
			Slug: slug,
		})

		if err != nil {
			c.JSON(500, gin.H{
				"error":  "Failed to save shortened URL",
				"reason": err.Error(),
			})
			return
		}

		// return shortened URL
		c.JSON(201, gin.H{
			"miniurl": buildShortenedURL(slug),
		})
	}
}

var ctx = context.Background()

func redirectURL(redisClient *redis.Client, postgresClient *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		urlQueryResult := &Slug{}

		url, err := redisClient.Get(ctx, slug).Result()
		if err == nil {
			fmt.Printf("Cache hit %s\n", slug)
			c.Redirect(301, url)
			return
		}
		if err != redis.Nil {
			fmt.Fprintf(os.Stderr, "Failed to retrieve URL from cache: %s\n", err.Error())
			c.JSON(500, gin.H{
				"error":  "Internal server error",
				"reason": err.Error(),
			})
			return
		}
		fmt.Printf("Cache miss %s\n", slug)

		err = postgresClient.Model(urlQueryResult).Column("url").Where("slug = ?", slug).Select()
		if err != nil {
			c.JSON(404, gin.H{
				"error":  "Page not found",
				"reason": "Please check the URL",
			})
			return
		}

		err = redisClient.Set(ctx, slug, urlQueryResult.Url, 0).Err()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to cache URL: %s\n", err.Error())
		}
		fmt.Printf("Cached %s\n", slug)

		c.Redirect(301, urlQueryResult.Url)
	}
}

func initRouter(redisClient *redis.Client, postgresClient *pg.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	router.POST("/mini", shortenURL(postgresClient))

	router.GET("/:slug", redirectURL(redisClient, postgresClient))

	return router
}
