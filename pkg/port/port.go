package port

import (
	"fmt"
	"net"
	"time"
)

// Config holds port connection configuration
type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
}

// Test checks if the given port is open
func Test(config Config) error {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	conn, err := net.DialTimeout("tcp", address, config.Timeout)
	
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	
	defer conn.Close()
	return nil
} 