package main

import (
	"fmt"
	"math/rand"
	"os"
)

func buildShortenedURL(slug string) string {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:8080"
	}
	return fmt.Sprintf("%s/%s", baseUrl, slug)
}

func generateSlug() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
