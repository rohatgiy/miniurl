package main

import (
	"net/url"

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

func redirectURL(postgresClient *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		urlQueryResult := &Slug{}
		err := postgresClient.Model(urlQueryResult).Column("url").Where("slug = ?", slug).Select()
		if err != nil {
			c.JSON(404, gin.H{
				"error":  "Page not found",
				"reason": "Please check the URL",
			})
			return
		}

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

	router.GET("/:slug", redirectURL(postgresClient))

	return router
}
