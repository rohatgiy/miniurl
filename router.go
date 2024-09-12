package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func rateLimiter(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(429, gin.H{
				"error":  "Too many requests",
				"reason": "Please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func shortenURL(postgresClient *pg.DB, req *ShortenURLRequestBody) (*ShortenURLResponse, *Error) {
	// validate URL
	url, err := url.ParseRequestURI(req.Url)
	if err != nil {
		return nil, &Error{
			Error:  "Invalid URL",
			Reason: err.Error(),
			Code:   400,
		}
	}

	// shorten URL/generate slug
	slug := generateSlug()

	// check if slug is already in use
	// use bloom filter?
	// add retry logic

	slugExists, err := checkIfSlugExists(postgresClient, slug)
	if slugExists {
		return nil, &Error{
			Error:  "Slug already in use",
			Reason: "Please try again",
			Code:   500,
		}
	}

	if err != pg.ErrNoRows {
		return nil, &Error{
			Error:  "Internal server error",
			Reason: err.Error(),
			Code:   500,
		}
	}

	// save to pgsql
	_, err = saveSlug(postgresClient, &Slug{
		Url:  url.String(),
		Slug: slug,
	})

	if err != nil {
		return nil, &Error{
			Error:  "Failed to save shortened URL",
			Reason: err.Error(),
			Code:   500,
		}
	}

	// return shortened URL
	return &ShortenURLResponse{
		MiniURL: buildShortenedURL(slug),
	}, nil
}

func bffShortenURLHandler(postgresClient *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body ShortenURLRequestBody
		parseError := c.ShouldBind(&body)

		if parseError != nil {
			c.HTML(400, "error.tmpl", gin.H{
				"error":  "Invalid request body",
				"reason": parseError.Error(),
			})
			return
		}

		response, err := shortenURL(postgresClient, &body)

		if err != nil {
			c.HTML(err.Code, "error.tmpl", gin.H{
				"error":  err.Error,
				"reason": err.Reason,
			})
			return
		}

		c.HTML(response.Code, "mini.tmpl", gin.H{
			"miniurl": response.MiniURL,
		})
	}
}

func shortenURLHandler(postgresClient *pg.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body ShortenURLRequestBody
		parseError := c.BindJSON(&body)

		if parseError != nil {
			c.JSON(400, gin.H{
				"error":  "Invalid request body",
				"reason": parseError.Error(),
			})
			return
		}

		response, err := shortenURL(postgresClient, &body)

		if err != nil {
			c.JSON(err.Code, gin.H{
				"error":  err.Error,
				"reason": err.Reason,
			})
			return
		}

		c.JSON(response.Code, gin.H{
			"miniurl": response.MiniURL,
		})
	}
}

func redirectURLHandler(redisClient *redis.Client, postgresClient *pg.DB) gin.HandlerFunc {
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
	limit := rate.Limit(100)
	limiter := rate.NewLimiter(limit, 10)

	router.LoadHTMLGlob("templates/*")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	router.Use(rateLimiter(limiter))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	router.GET("/mini", func(c *gin.Context) {
		c.Redirect(301, "/")
	})

	router.POST("/mini", bffShortenURLHandler(postgresClient))

	router.POST("/api/mini", shortenURLHandler(postgresClient))

	router.GET("/:slug", redirectURLHandler(redisClient, postgresClient))

	return router
}
