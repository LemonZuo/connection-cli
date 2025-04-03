package redis

import (
	"context"
	"fmt"
	"time"

	"connection-cli/pkg/logger"
	"github.com/go-redis/redis/v8"
)

// Config holds Redis connection configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
	Timeout  time.Duration
}

// Test checks if the Redis server can be connected
func Test(config Config) error {
	// 初始化日志
	logger.Init()

	logger.Info("Testing Redis connection to %s:%d with username: %s, DB: %d", config.Host, config.Port, config.Username, config.DB)
	
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
		logger.Debug("Using default timeout: %s", config.Timeout)
	}

	// 打印连接参数（密码不打印）
	logger.Debug("Redis connection details - Host: %s, Port: %d, Username: %s, DB: %d, Timeout: %s", 
		config.Host, config.Port, config.Username, config.DB, config.Timeout)

	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username:    config.Username,
		Password:    config.Password,
		DB:          config.DB,
		DialTimeout: config.Timeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	logger.Info("Attempting to ping Redis server...")
	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Error("Failed to connect to Redis: %v", err)
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Successfully connected to Redis")
	return nil
} 