package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl  string
	DatabaseName string
	Port         string
	JWTSecret    string
}

func NewConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	return &Config{
		DatabaseUrl:  os.Getenv("DATABASE_URL"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		Port:         os.Getenv("PORT"),
	}
}
