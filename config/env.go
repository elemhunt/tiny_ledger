package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	/*
		Constructor to load environment variables from a .env file
	*/
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func GetEnv(key, fallback string) string {
	/*
		Constructor to get environment variable by key else giving a fallback default for
		that variable.
		Return: Variable of type string
	*/
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}
