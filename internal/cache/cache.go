package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type (
	Cache struct {
		redisClient *redis.Client
	}
)

func NewCache(redisClient *redis.Client) *Cache {
	return &Cache{
		redisClient: redisClient,
	}
}

func (cache *Cache) PushToCache(ctx context.Context, orderUidStr string, data []byte) error {
	if err := cache.redisClient.
		Set(ctx, orderUidStr, data, 0).
		Err(); err != nil {
		return fmt.Errorf("failed set data to Redis: ")
	}
	return nil
}

func (cache *Cache) PullFromCache(ctx context.Context, orderUidStr string) ([]byte, error) {
	data, err := cache.redisClient.
		Get(ctx, orderUidStr).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get data from Redis: ")
	}
	return data, nil
}
