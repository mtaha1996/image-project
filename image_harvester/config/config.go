package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleAPIKey       string
	SearchEngineID     string
	BaseURL            string
	PostgresDB         string
	PostgresUser       string
	PostgresPassword   string
	PostgresHost       string
	PostgresPort       string
	ImageDirectoryName string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return Config{
		GoogleAPIKey:       os.Getenv("GOOGLE_API_KEY"),
		SearchEngineID:     os.Getenv("SEARCH_ENGINE_ID"),
		BaseURL:            os.Getenv("BASE_URL"),
		PostgresDB:         os.Getenv("IMAGE_POSTGRES_DB"),
		PostgresUser:       os.Getenv("IMAGE_POSTGRES_USER"),
		PostgresPassword:   os.Getenv("IMAGE_POSTGRES_PASSWORD"),
		PostgresHost:       os.Getenv("IMAGE_POSTGRES_HOST"),
		PostgresPort:       os.Getenv("IMAGE_POSTGRES_PORT"),
		ImageDirectoryName: os.Getenv("IMAGE_DIRECTORY_NAME"),
	}
}
