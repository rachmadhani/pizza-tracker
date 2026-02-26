package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	Port       string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

var AppConfig *Config

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load environment variables")
	}

	AppConfig = &Config{
		DBHost:     GetEnv("DB_HOST", "localhost"),
		Port:       GetEnv("APP_PORT", "8080"),
		DBPort:     GetEnv("DB_PORT", "3306"),
		DBName:     GetEnv("DB_NAME", "pizza_tracker"),
		DBUser:     GetEnv("DB_USER", "root"),
		DBPassword: GetEnv("DB_PASSWORD", ""),
	}
	return AppConfig
}

func GetEnv(key, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}
