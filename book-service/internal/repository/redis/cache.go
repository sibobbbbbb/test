package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sibobbbbbb/backend-engineer-challenge/book-service/internal/domain"
)

type bookCache struct {
	client *redis.Client
}

// NewBookCache creates a new instance of book cache
func NewBookCache(redisURL string, redisPassword string, redisDB int) (*bookCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test koneksi
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &bookCache{
		client: client,
	}, nil
}

// Set menyimpan buku di cache
func (c *bookCache) Set(ctx context.Context, id string, book *domain.Book, expiration time.Duration) error {
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("failed to marshal book: %w", err)
	}

	key := fmt.Sprintf("book:%s", id)
	return c.client.Set(ctx, key, bookJSON, expiration).Err()
}

// Get mengambil buku dari cache
func (c *bookCache) Get(ctx context.Context, id string) (*domain.Book, error) {
	key := fmt.Sprintf("book:%s", id)
	bookJSON, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("book not found in cache")
		}
		return nil, fmt.Errorf("failed to get book from cache: %w", err)
	}

	book := &domain.Book{}
	if err := json.Unmarshal([]byte(bookJSON), book); err != nil {
		return nil, fmt.Errorf("failed to unmarshal book: %w", err)
	}

	return book, nil
}

// Delete menghapus buku dari cache
func (c *bookCache) Delete(ctx context.Context, id string) error {
	key := fmt.Sprintf("book:%s", id)
	return c.client.Del(ctx, key).Err()
}

// SetSearchResults menyimpan hasil pencarian di cache
func (c *bookCache) SetSearchResults(ctx context.Context, query string, page int, limit int, books []*domain.Book, total int, expiration time.Duration) error {
	result := struct {
		Books []*domain.Book `json:"books"`
		Total int            `json:"total"`
	}{
		Books: books,
		Total: total,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal search results: %w", err)
	}

	key := fmt.Sprintf("search:%s:page:%d:limit:%d", query, page, limit)
	return c.client.Set(ctx, key, resultJSON, expiration).Err()
}

// GetSearchResults mengambil hasil pencarian dari cache
func (c *bookCache) GetSearchResults(ctx context.Context, query string, page int, limit int) ([]*domain.Book, int, error) {
	key := fmt.Sprintf("search:%s:page:%d:limit:%d", query, page, limit)
	resultJSON, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, 0, fmt.Errorf("search results not found in cache")
		}
		return nil, 0, fmt.Errorf("failed to get search results from cache: %w", err)
	}

	result := struct {
		Books []*domain.Book `json:"books"`
		Total int            `json:"total"`
	}{}

	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal search results: %w", err)
	}

	return result.Books, result.Total, nil
}

// SetBooksList menyimpan daftar buku di cache
func (c *bookCache) SetBooksList(ctx context.Context, page int, limit int, books []*domain.Book, total int, expiration time.Duration) error {
	result := struct {
		Books []*domain.Book `json:"books"`
		Total int            `json:"total"`
	}{
		Books: books,
		Total: total,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal books list: %w", err)
	}

	key := fmt.Sprintf("books:page:%d:limit:%d", page, limit)
	return c.client.Set(ctx, key, resultJSON, expiration).Err()
}

// GetBooksList mengambil daftar buku dari cache
func (c *bookCache) GetBooksList(ctx context.Context, page int, limit int) ([]*domain.Book, int, error) {
	key := fmt.Sprintf("books:page:%d:limit:%d", page, limit)
	resultJSON, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, 0, fmt.Errorf("books list not found in cache")
		}
		return nil, 0, fmt.Errorf("failed to get books list from cache: %w", err)
	}

	result := struct {
		Books []*domain.Book `json:"books"`
		Total int            `json:"total"`
	}{}

	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal books list: %w", err)
	}

	return result.Books, result.Total, nil
}

// InvalidateCache menghapus semua cache terkait dengan buku
func (c *bookCache) InvalidateCache(ctx context.Context) error {
	// Ini adalah contoh implementasi sederhana
	// Dalam produksi, mungkin perlu pendekatan yang lebih halus
	pattern := "book:*"
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys: %w", err)
		}
	}

	// Hapus juga cache daftar dan pencarian
	listPattern := "books:page:*"
	searchPattern := "search:*"

	listKeys, err := c.client.Keys(ctx, listPattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get list keys: %w", err)
	}

	if len(listKeys) > 0 {
		if err := c.client.Del(ctx, listKeys...).Err(); err != nil {
			return fmt.Errorf("failed to delete list keys: %w", err)
		}
	}

	searchKeys, err := c.client.Keys(ctx, searchPattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get search keys: %w", err)
	}

	if len(searchKeys) > 0 {
		if err := c.client.Del(ctx, searchKeys...).Err(); err != nil {
			return fmt.Errorf("failed to delete search keys: %w", err)
		}
	}

	return nil
}

// Close menutup koneksi Redis
func (c *bookCache) Close() error {
	return c.client.Close()
}