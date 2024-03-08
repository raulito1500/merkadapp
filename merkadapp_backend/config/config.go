package server

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Usar las config dentro de las aplicaciones
type Config struct {
	DatabaseUrl string
	Port        string
	JWTSecret   string
}

func NewConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading enviroments")
	}
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	return &Config{
		DatabaseUrl: DATABASE_URL,
		JWTSecret:   JWT_SECRET,
	}
}
