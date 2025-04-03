package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Config holds HTTP connection configuration
type Config struct {
	URL     string
	Method  string
	Timeout time.Duration
}

// Test checks if the given HTTP URL is accessible
func Test(config Config) error {
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	if config.Method == "" {
		config.Method = http.MethodGet
	}

	client := &http.Client{
		Timeout: config.Timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", config.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received HTTP status code %d from %s", resp.StatusCode, config.URL)
	}

	return nil
} 