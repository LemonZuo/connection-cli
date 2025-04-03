package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=%d",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
		config.SSLMode,
		int(config.Timeout.Seconds()),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to create PostgreSQL connection: %w", err)
	}
	defer db.Close()

	// Set connection timeout
	db.SetConnMaxLifetime(config.Timeout)

	// Ping to test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return nil
} 