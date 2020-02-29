package cache

import (
	"fmt"

	"github.com/go-redis/redis"

	"github.com/superhero-chat/internal/config"
)

// Cache is the Redis client.
type Cache struct {
	Redis               *redis.Client
	OnlineUserKeyFormat string
}

// NewCache creates a client connection to Redis.
func NewCache(cfg *config.Config) (cache *Cache, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s%s", cfg.Cache.Address, cfg.Cache.Port),
		Password:     cfg.Cache.Password,
		DB:           cfg.Cache.DB,
		PoolSize:     cfg.Cache.PoolSize,
		MinIdleConns: cfg.Cache.MinimumIdleConnections,
		MaxRetries:   cfg.Cache.MaximumRetries,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Cache{
		Redis:               client,
		OnlineUserKeyFormat: cfg.Cache.OnlineUserKeyFormat,
	}, nil
}
