// user-service/internal/repository/redis/token_repo.go
package redis

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/domain"
)

type redisRepository struct {
    client *redis.Client
}

// NewRedisRepository membuat instance baru Redis repository
func NewRedisRepository(redisURL, redisPassword string, redisDB int) domain.UserRepository {
    client := redis.NewClient(&redis.Options{
        Addr:     redisURL,
        Password: redisPassword,
        DB:       redisDB,
    })

    return &redisRepository{
        client: client,
    }
}

// StoreToken menyimpan token di Redis
func (r *redisRepository) StoreToken(ctx context.Context, userID, token string, expiry time.Duration) error {
    key := fmt.Sprintf("token:%s", token)
    return r.client.Set(ctx, key, userID, expiry).Err()
}

// GetUserIDByToken mendapatkan userID dari token
func (r *redisRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
    key := fmt.Sprintf("token:%s", token)
    userID, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return "", fmt.Errorf("token not found or expired")
        }
        return "", fmt.Errorf("failed to get user ID by token: %w", err)
    }
    return userID, nil
}

// DeleteToken menghapus token dari Redis
func (r *redisRepository) DeleteToken(ctx context.Context, token string) error {
    key := fmt.Sprintf("token:%s", token)
    return r.client.Del(ctx, key).Err()
}