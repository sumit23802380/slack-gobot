package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetupConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
