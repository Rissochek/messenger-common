package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvFile() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("no .env file found")
	}
}

func GetKeyFromEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%v value is not set in .env file.", key)
	}
	return value
}