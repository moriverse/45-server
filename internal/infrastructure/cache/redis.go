package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/moriverse/45-server/internal/infrastructure/config"
)

// NewRedisClient creates and returns a new Redis client based on the provided configuration.
func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return rdb
}
