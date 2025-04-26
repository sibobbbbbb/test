package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	// Server
	ServerAddress string
	ServerPort    string
	GRPCPort      string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Cache
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// Category Service
	CategoryServiceHost string
	CategoryServicePort string

	// JWT
	JWTSecret    string
	JWTExpiry    time.Duration
	RefreshToken time.Duration
}

// LoadConfig memuat konfigurasi dari environment variables
func LoadConfig() *Config {
	// Jika environment variable ada, gunakan nilai tersebut
	// Jika tidak, gunakan nilai default
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY", "3600"))
	refreshToken, _ := strconv.Atoi(getEnv("REFRESH_TOKEN", "86400"))

	return &Config{
		// Server
		ServerAddress: getEnv("SERVER_ADDRESS", "0.0.0.0"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		GRPCPort:      getEnv("GRPC_PORT", "50051"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "book_service"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Cache
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		// Category Service
		CategoryServiceHost: getEnv("CATEGORY_SERVICE_HOST", "localhost"),
		CategoryServicePort: getEnv("CATEGORY_SERVICE_PORT", "50052"),

		// JWT
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		JWTExpiry:    time.Duration(jwtExpiry) * time.Second,
		RefreshToken: time.Duration(refreshToken) * time.Second,
	}
}

// GetDBConnectionString mengembalikan string koneksi database
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// GetRedisConnectionString mengembalikan string koneksi Redis
func (c *Config) GetRedisConnectionString() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// GetCategoryServiceAddress mengembalikan alamat category service
func (c *Config) GetCategoryServiceAddress() string {
	return fmt.Sprintf("%s:%s", c.CategoryServiceHost, c.CategoryServicePort)
}

// getEnv mengambil nilai environment variable atau default value jika tidak ada
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}