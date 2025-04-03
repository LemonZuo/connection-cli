package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connection-cli/pkg/logger"
)

// Config holds HTTP connection configuration
type Config struct {
	URL     string
	Method  string
	Timeout time.Duration
}

// Test checks if the given HTTP URL is accessible
func Test(config Config) error {
	// 初始化日志
	logger.Init()
	
	logger.Info("Testing HTTP connection to %s using method %s", config.URL, config.Method)
	
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
		logger.Debug("Using default timeout: %s", config.Timeout)
	}

	if config.Method == "" {
		config.Method = http.MethodGet
		logger.Debug("Using default HTTP method: %s", config.Method)
	}

	logger.Debug("HTTP connection details - URL: %s, Method: %s, Timeout: %s", 
		config.URL, config.Method, config.Timeout)

	client := &http.Client{
		Timeout: config.Timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	logger.Info("Creating new HTTP request...")
	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, nil)
	if err != nil {
		logger.Error("Failed to create HTTP request: %v", err)
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	logger.Info("Sending HTTP request to %s...", config.URL)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Failed to connect to %s: %v", config.URL, err)
		return fmt.Errorf("failed to connect to %s: %w", config.URL, err)
	}
	defer resp.Body.Close()

	logger.Info("Received HTTP status code %d from %s", resp.StatusCode, config.URL)
	if resp.StatusCode >= 400 {
		logger.Error("HTTP request failed with status code %d", resp.StatusCode)
		return fmt.Errorf("received HTTP status code %d from %s", resp.StatusCode, config.URL)
	}

	logger.Info("Successfully connected to %s", config.URL)
	return nil
} 