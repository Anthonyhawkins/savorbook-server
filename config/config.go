package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func Get(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Unable to load .env file", err)
	}
	return os.Getenv(key)
}
