package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Server
	ServerAddress string
	ServerPort    string

	// Services
	BookServiceHost     string
	BookServicePort     string
	CategoryServiceHost string
	CategoryServicePort string
	UserServiceHost     string
	UserServicePort     string

	// Timeout
	ContextTimeout time.Duration
}

func LoadConfig() *Config {
	timeout, _ := strconv.Atoi(getEnv("CONTEXT_TIMEOUT", "30"))

	return &Config{
		// Server
		ServerAddress: getEnv("SERVER_ADDRESS", "0.0.0.0"),
		ServerPort:    getEnv("SERVER_PORT", "8000"),

		// Services
		BookServiceHost:     getEnv("BOOK_SERVICE_HOST", "localhost"),
		BookServicePort:     getEnv("BOOK_SERVICE_PORT", "50051"),
		CategoryServiceHost: getEnv("CATEGORY_SERVICE_HOST", "localhost"),
		CategoryServicePort: getEnv("CATEGORY_SERVICE_PORT", "50052"),
		UserServiceHost:     getEnv("USER_SERVICE_HOST", "localhost"),
		UserServicePort:     getEnv("USER_SERVICE_PORT", "50053"),

		// Timeout
		ContextTimeout: time.Duration(timeout) * time.Second,
	}
}

func (c *Config) GetBookServiceAddress() string {
	return c.BookServiceHost + ":" + c.BookServicePort
}

func (c *Config) GetCategoryServiceAddress() string {
	return c.CategoryServiceHost + ":" + c.CategoryServicePort
}

func (c *Config) GetUserServiceAddress() string {
	return c.UserServiceHost + ":" + c.UserServicePort
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}