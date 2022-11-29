package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHostUrl string
	ImagesDir     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
		return nil, err
	}

	log.Println("Loaded .env file")

	return &Config{
		ServerHostUrl: os.Getenv("SERVER_HOST_URL"),
		ImagesDir:     os.Getenv("IMAGES_DIR"),
	}, nil
}
