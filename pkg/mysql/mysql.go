package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s", 
		config.Username, 
		config.Password, 
		config.Host, 
		config.Port, 
		config.Database,
		config.Timeout.String(),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to create MySQL connection: %w", err)
	}
	defer db.Close()

	// Set connection timeout
	db.SetConnMaxLifetime(config.Timeout)

	// Ping to test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	return nil
} 