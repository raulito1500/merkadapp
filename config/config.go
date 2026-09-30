package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl         string
	DatabaseName        string
	Port                string
	FirebaseProjectID   string
	FirebaseClientEmail string
	FirebasePrivateKey  string
}

func NewConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	return &Config{
		DatabaseUrl:         os.Getenv("DATABASE_URL"),
		DatabaseName:        os.Getenv("DATABASE_NAME"),
		Port:                os.Getenv("PORT"),
		FirebaseProjectID:   os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseClientEmail: os.Getenv("FIREBASE_CLIENT_EMAIL"),
		FirebasePrivateKey:  os.Getenv("FIREBASE_PRIVATE_KEY"),
	}
}
