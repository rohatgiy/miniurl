package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env file: %s\n", err.Error())
		os.Exit(1)
		return
	}
	redis := initRedis()
	postgres := initPostgres()
	router := initRouter(redis, postgres)

	router.Run()
}
