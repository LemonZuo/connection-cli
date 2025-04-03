package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"connection-cli/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
)

// Config holds MySQL connection configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Timeout  time.Duration
}

// Test checks if the MySQL server can be connected
func Test(config Config) error {
	// 初始化日志
	logger.Init()
	
	logger.Info("Testing MySQL connection to %s:%d", config.Host, config.Port)
	
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
		logger.Debug("Using default timeout: %s", config.Timeout)
	}

	logger.Debug("MySQL connection details - Host: %s, Port: %d, Username: %s, Database: %s, Timeout: %s",
		config.Host, config.Port, config.Username, config.Database, config.Timeout)

	// 不包括密码信息，避免安全风险
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s", 
		config.Username, 
		config.Password, 
		config.Host, 
		config.Port, 
		config.Database,
		config.Timeout.String(),
	)

	logger.Info("Opening MySQL connection...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to create MySQL connection: %v", err)
		return fmt.Errorf("failed to create MySQL connection: %w", err)
	}
	defer db.Close()

	// Set connection timeout
	db.SetConnMaxLifetime(config.Timeout)
	logger.Debug("Connection max lifetime set to: %s", config.Timeout)

	// Ping to test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	logger.Info("Pinging MySQL server...")
	if err := db.PingContext(ctx); err != nil {
		logger.Error("Failed to connect to MySQL: %v", err)
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	logger.Info("Successfully connected to MySQL")
	return nil
} 