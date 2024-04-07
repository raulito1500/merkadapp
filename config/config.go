package config

import (
	"os"
)

type Config struct {
	DatabaseUrl  string
	DatabaseName string
	Port         string
	JWTSecret    string
}

func NewConfig() *Config {
	return &Config{
		DatabaseUrl:  os.Getenv("DATABASE_URL"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		Port:         os.Getenv("PORT"),
	}
}
