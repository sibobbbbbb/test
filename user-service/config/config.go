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

    // Redis
    RedisHost     string
    RedisPort     string
    RedisPassword string
    RedisDB       int

    // JWT
    JWTSecret    string
    JWTExpiry    time.Duration
    RefreshToken time.Duration

    // Context
    ContextTimeout time.Duration
}

func LoadConfig() *Config {
    redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
    timeout, _ := strconv.Atoi(getEnv("CONTEXT_TIMEOUT", "30"))
    jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY", "3600"))
    refreshToken, _ := strconv.Atoi(getEnv("REFRESH_TOKEN", "86400"))

    return &Config{
        // Server
        ServerAddress: getEnv("SERVER_ADDRESS", "0.0.0.0"),
        ServerPort:    getEnv("SERVER_PORT", "8082"),
        GRPCPort:      getEnv("GRPC_PORT", "50053"),

        // Database
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "user_service"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

        // Redis
        RedisHost:     getEnv("REDIS_HOST", "localhost"),
        RedisPort:     getEnv("REDIS_PORT", "6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       redisDB,

        // JWT
        JWTSecret:    getEnv("JWT_SECRET", "secret"),
        JWTExpiry:    time.Duration(jwtExpiry) * time.Second,
        RefreshToken: time.Duration(refreshToken) * time.Second,

        // Context
        ContextTimeout: time.Duration(timeout) * time.Second,
    }
}

func (c *Config) GetDBConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
    )
}

func (c *Config) GetRedisConnectionString() string {
    return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}