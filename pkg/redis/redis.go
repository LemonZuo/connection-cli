package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config holds Redis connection configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	Timeout  time.Duration
}

// Test checks if the Redis server can be connected
func Test(config Config) error {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:    config.Password,
		DB:          config.DB,
		DialTimeout: config.Timeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
} 