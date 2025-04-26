package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

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

    // Context
    ContextTimeout time.Duration

    // JWT
    JWTSecret    string
    JWTExpiry    time.Duration
}

func LoadConfig() *Config {
    timeout, _ := strconv.Atoi(getEnv("CONTEXT_TIMEOUT", "30"))
    jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY", "3600"))

    return &Config{
        // Server
        ServerAddress: getEnv("SERVER_ADDRESS", "0.0.0.0"),
        ServerPort:    getEnv("SERVER_PORT", "8081"),
        GRPCPort:      getEnv("GRPC_PORT", "50052"),

        // Database
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "category_service"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

        // Context
        ContextTimeout: time.Duration(timeout) * time.Second,

        // JWT
        JWTSecret:    getEnv("JWT_SECRET", "secret"),
        JWTExpiry:    time.Duration(jwtExpiry) * time.Second,
    }
}

func (c *Config) GetDBConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
    )
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}