package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/sibobbbbbb/backend-engineer-challenge/category-service/internal/domain"
)

type categoryCache struct {
    client *redis.Client
}

func NewCategoryCache(redisURL string, redisPassword string, redisDB int) (*categoryCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     redisURL,
        Password: redisPassword,
        DB:       redisDB,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return &categoryCache{
        client: client,
    }, nil
}

func (c *categoryCache) Set(ctx context.Context, id string, category *domain.Category, expiration time.Duration) error {
    categoryJSON, err := json.Marshal(category)
    if err != nil {
        return fmt.Errorf("failed to marshal category: %w", err)
    }

    key := fmt.Sprintf("category:%s", id)
    return c.client.Set(ctx, key, categoryJSON, expiration).Err()
}

func (c *categoryCache) Get(ctx context.Context, id string) (*domain.Category, error) {
    key := fmt.Sprintf("category:%s", id)
    categoryJSON, err := c.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, fmt.Errorf("category not found in cache")
        }
        return nil, fmt.Errorf("failed to get category from cache: %w", err)
    }

    category := &domain.Category{}
    if err := json.Unmarshal([]byte(categoryJSON), category); err != nil {
        return nil, fmt.Errorf("failed to unmarshal category: %w", err)
    }

    return category, nil
}

func (c *categoryCache) Delete(ctx context.Context, id string) error {
    key := fmt.Sprintf("category:%s", id)
    return c.client.Del(ctx, key).Err()
}

func (c *categoryCache) SetCategoriesList(ctx context.Context, page int, limit int, categories []*domain.Category, total int, expiration time.Duration) error {
    result := struct {
        Categories []*domain.Category `json:"categories"`
        Total int `json:"total"`
    }{
        Categories: categories,
        Total: total,
    }

    resultJSON, err := json.Marshal(result)
    if err != nil {
        return fmt.Errorf("failed to marshal categories list: %w", err)
    }

    key := fmt.Sprintf("categories:page:%d:limit:%d", page, limit)
    return c.client.Set(ctx, key, resultJSON, expiration).Err()
}

func (c *categoryCache) GetCategoriesList(ctx context.Context, page int, limit int) ([]*domain.Category, int, error) {
    key := fmt.Sprintf("categories:page:%d:limit:%d", page, limit)
    resultJSON, err := c.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return nil, 0, fmt.Errorf("categories list not found in cache")
        }
        return nil, 0, fmt.Errorf("failed to get categories list from cache: %w", err)
    }

    result := struct {
        Categories []*domain.Category `json:"categories"`
        Total int `json:"total"`
    }{}

    if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
        return nil, 0, fmt.Errorf("failed to unmarshal categories list: %w", err)
    }

    return result.Categories, result.Total, nil
}

func (c *categoryCache) InvalidateCache(ctx context.Context) error {
    pattern := "category:*"
    keys, err := c.client.Keys(ctx, pattern).Result()
    if err != nil {
        return fmt.Errorf("failed to get keys: %w", err)
    }

    if len(keys) > 0 {
        if err := c.client.Del(ctx, keys...).Err(); err != nil {
            return fmt.Errorf("failed to delete keys: %w", err)
        }
    }

    listPattern := "categories:page:*"
    listKeys, err := c.client.Keys(ctx, listPattern).Result()
    if err != nil {
        return fmt.Errorf("failed to get list keys: %w", err)
    }

    if len(listKeys) > 0 {
        if err := c.client.Del(ctx, listKeys...).Err(); err != nil {
            return fmt.Errorf("failed to delete list keys: %w", err)
        }
    }

    return nil
}

func (c *categoryCache) Close() error {
    return c.client.Close()
}