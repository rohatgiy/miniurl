package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env file: %s\n", err.Error())
		os.Exit(1)
		return
	}

	env := os.Getenv("ENV")

	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	redis := initRedis()
	postgres := initPostgres()
	router := initRouter(redis, postgres)
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"192.168.0.0/20"})

	router.Run()
}
