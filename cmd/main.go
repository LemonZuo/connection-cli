package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"connection-cli/pkg/http"
	"connection-cli/pkg/logger"
	"connection-cli/pkg/mysql"
	"connection-cli/pkg/port"
	"connection-cli/pkg/postgres"
	"connection-cli/pkg/redis"
)

// Version is the application version, set at build time
var Version = "dev"

// Config holds the application configuration
type Config struct {
	Mode        string
	Host        string
	Port        int
	Username    string
	Password    string
	Database    string
	URL         string
	Timeout     time.Duration
	SSLMode     string
	RedisDB     int
	HTTPMethod  string
	ShowVersion bool
}

func main() {
	// 初始化日志
	logger.Init()
	defer logger.Close()

	var config Config

	// Define command-line flags
	flag.StringVar(&config.Mode, "mode", "", "Testing mode: mysql, postgres, redis, port, http")
	flag.StringVar(&config.Host, "host", "localhost", "Host to connect to")
	flag.IntVar(&config.Port, "port", 0, "Port to connect to")
	flag.StringVar(&config.Username, "username", "", "Username for database connections")
	flag.StringVar(&config.Password, "password", "", "Password for database connections")
	flag.StringVar(&config.Database, "database", "", "Database name for database connections")
	flag.StringVar(&config.URL, "url", "", "URL for HTTP testing")
	flag.DurationVar(&config.Timeout, "timeout", 5*time.Second, "Connection timeout")
	flag.StringVar(&config.SSLMode, "sslmode", "disable", "SSL mode for PostgreSQL")
	flag.IntVar(&config.RedisDB, "redis-db", 0, "Redis database number")
	flag.StringVar(&config.HTTPMethod, "http-method", "GET", "HTTP method for HTTP testing")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version information and exit")

	// Parse from environment variables if present
	flag.Parse()

	// Show version if requested
	if config.ShowVersion {
		logger.Info("Connection CLI version %s", Version)
		fmt.Printf("Connection CLI version %s\n", Version)
		os.Exit(0)
	}

	// Override with environment variables if set
	if os.Getenv("MODE") != "" {
		config.Mode = os.Getenv("MODE")
	}
	if os.Getenv("HOST") != "" {
		config.Host = os.Getenv("HOST")
	}
	if os.Getenv("PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err == nil {
			config.Port = port
		}
	}
	if os.Getenv("USERNAME") != "" {
		config.Username = os.Getenv("USERNAME")
	}
	if os.Getenv("PASSWORD") != "" {
		config.Password = os.Getenv("PASSWORD")
	}
	if os.Getenv("DATABASE") != "" {
		config.Database = os.Getenv("DATABASE")
	}
	if os.Getenv("URL") != "" {
		config.URL = os.Getenv("URL")
	}
	if os.Getenv("TIMEOUT") != "" {
		timeout, err := time.ParseDuration(os.Getenv("TIMEOUT"))
		if err == nil {
			config.Timeout = timeout
		}
	}
	if os.Getenv("SSLMODE") != "" {
		config.SSLMode = os.Getenv("SSLMODE")
	}
	if os.Getenv("REDIS_DB") != "" {
		redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err == nil {
			config.RedisDB = redisDB
		}
	}
	if os.Getenv("HTTP_METHOD") != "" {
		config.HTTPMethod = os.Getenv("HTTP_METHOD")
	}

	// If mode is empty, display help and exit
	if config.Mode == "" || config.Mode == "help" {
		logger.Info("Showing help information")
		fmt.Printf("Connection CLI %s - A tool for testing connectivity to various services\n\n", Version)
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println("\nExample: connection-cli -mode=port -host=example.com -port=80")
		os.Exit(0)
	}

	logger.Info("Starting connection test: mode=%s, host=%s, port=%d", config.Mode, config.Host, config.Port)

	// Validate required parameters
	if config.Mode != "http" && config.Port == 0 {
		logger.Error("Port is required for mode: %s", config.Mode)
		fmt.Fprintf(os.Stderr, "Error: Port is required for mode: %s\n", config.Mode)
		os.Exit(1)
	}

	// Validate URL for HTTP mode
	if config.Mode == "http" && config.URL == "" {
		logger.Error("URL is required for HTTP mode")
		fmt.Fprintf(os.Stderr, "Error: URL is required for HTTP mode\n")
		os.Exit(1)
	}

	// Health check based on the selected mode
	var testErr error
	switch strings.ToLower(config.Mode) {
	case "mysql":
		logger.Info("Testing MySQL connection to %s:%d", config.Host, config.Port)
		testErr = mysql.Test(mysql.Config{
			Host:     config.Host,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
			Database: config.Database,
			Timeout:  config.Timeout,
		})
	case "postgres":
		logger.Info("Testing PostgreSQL connection to %s:%d", config.Host, config.Port)
		testErr = postgres.Test(postgres.Config{
			Host:     config.Host,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
			Database: config.Database,
			SSLMode:  config.SSLMode,
			Timeout:  config.Timeout,
		})
	case "redis":
		logger.Info("Testing Redis connection to %s:%d", config.Host, config.Port)
		testErr = redis.Test(redis.Config{
			Host:     config.Host,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
			DB:       config.RedisDB,
			Timeout:  config.Timeout,
		})
	case "port":
		logger.Info("Testing port connectivity to %s:%d", config.Host, config.Port)
		testErr = port.Test(port.Config{
			Host:    config.Host,
			Port:    config.Port,
			Timeout: config.Timeout,
		})
	case "http":
		logger.Info("Testing HTTP connection to %s", config.URL)
		testErr = http.Test(http.Config{
			URL:     config.URL,
			Method:  config.HTTPMethod,
			Timeout: config.Timeout,
		})
	default:
		logger.Error("Unsupported mode: %s", config.Mode)
		fmt.Fprintf(os.Stderr, "Error: Unsupported mode: %s\n", config.Mode)
		os.Exit(1)
	}

	if testErr != nil {
		logger.Error("Connection test failed: %v", testErr)
		fmt.Fprintf(os.Stderr, "Error: Connection test failed: %v\n", testErr)
		os.Exit(1)
	}

	successMsg := fmt.Sprintf("Successfully connected to %s", config.Mode)
	logger.Info(successMsg)
	fmt.Println(successMsg)
	os.Exit(0)
}
