package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName     string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	BotToken   string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		BotToken:   os.Getenv("BOT_TOKEN"),
	}
}
