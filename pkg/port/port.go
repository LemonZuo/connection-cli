package port

import (
	"fmt"
	"net"
	"time"

	"connection-cli/pkg/logger"
)

// Config holds port connection configuration
type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
}

// Test checks if the given port is open
func Test(config Config) error {
	// 初始化日志
	logger.Init()
	
	logger.Info("Testing TCP port connectivity to %s:%d", config.Host, config.Port)
	
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
		logger.Debug("Using default timeout: %s", config.Timeout)
	}

	logger.Debug("Port connection details - Host: %s, Port: %d, Timeout: %s", 
		config.Host, config.Port, config.Timeout)

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	logger.Info("Attempting to connect to %s...", address)
	
	conn, err := net.DialTimeout("tcp", address, config.Timeout)
	if err != nil {
		logger.Error("Failed to connect to %s: %v", address, err)
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	
	defer conn.Close()
	logger.Info("Successfully connected to port %s:%d", config.Host, config.Port)
	return nil
} 