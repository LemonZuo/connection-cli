package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"connection-cli/pkg/logger"
	_ "github.com/lib/pq"
)

// Config holds PostgreSQL connection configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  string
	Timeout  time.Duration
}

// Test checks if the PostgreSQL server can be connected
func Test(config Config) error {
	// 初始化日志
	logger.Init()
	
	logger.Info("Testing PostgreSQL connection to %s:%d", config.Host, config.Port)
	
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
		logger.Debug("Using default timeout: %s", config.Timeout)
	}

	if config.SSLMode == "" {
		config.SSLMode = "disable"
		logger.Debug("Using default SSL mode: %s", config.SSLMode)
	}

	logger.Debug("PostgreSQL connection details - Host: %s, Port: %d, Username: %s, Database: %s, SSLMode: %s, Timeout: %s",
		config.Host, config.Port, config.Username, config.Database, config.SSLMode, config.Timeout)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		config.SSLMode,
		int(config.Timeout.Seconds()),
	)

	logger.Info("Opening PostgreSQL connection...")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to create PostgreSQL connection: %v", err)
		return fmt.Errorf("failed to create PostgreSQL connection: %w", err)
	}
	defer db.Close()

	// Set connection timeout
	db.SetConnMaxLifetime(config.Timeout)
	logger.Debug("Connection max lifetime set to: %s", config.Timeout)

	// Ping to test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	logger.Info("Pinging PostgreSQL server...")
	if err := db.PingContext(ctx); err != nil {
		logger.Error("Failed to connect to PostgreSQL: %v", err)
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL")
	return nil
} 