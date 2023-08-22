package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName          string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	BotToken        string
	MaxOpenConns    string
	MaxIdleConns    string
	ConnMaxIdleTime string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("after loading config")
	return Config{
		DBName:          os.Getenv("DB_NAME"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		BotToken:        os.Getenv("BOT_TOKEN"),
		MaxOpenConns:    os.Getenv("MAX_OPEN_CONNS"),
		MaxIdleConns:    os.Getenv("MAX_IDLE_CONNS"),
		ConnMaxIdleTime: os.Getenv("CONN_MAX_IDLE_TIME"),
	}
}
