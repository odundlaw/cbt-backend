// Package env for managing all environment variables
package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
}

func GetString[T any](key string, fallback T) T {
	loadEnv()

	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	switch any(fallback).(type) {
	case string:
		return any(value).(T)
	case int:
		if i, err := strconv.Atoi(value); err == nil {
			return any(i).(T)
		}
	case bool:
		if b, err := strconv.ParseBool(value); err == nil {
			return any(b).(T)
		}
	}

	return fallback
}
