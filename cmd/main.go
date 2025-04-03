package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"connection-cli/pkg/http"
	"connection-cli/pkg/mysql"
	"connection-cli/pkg/port"
	"connection-cli/pkg/postgres"
	"connection-cli/pkg/redis"
)

// Version is the application version, set at build time
var Version = "dev"

// Config holds the application configuration
type Config struct {
	Mode       string
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	URL        string
	Timeout    time.Duration
	SSLMode    string
	RedisDB    int
	HTTPMethod string
	ShowVersion bool
}

func main() {
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
		fmt.Printf("Connection CLI %s - A tool for testing connectivity to various services\n\n", Version)
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println("\nExample: connection-cli -mode=port -host=example.com -port=80")
		os.Exit(0)
	}

	// Validate required parameters
	if config.Mode != "http" && config.Port == 0 {
		log.Fatal("Port is required for mode: ", config.Mode)
	}

	// Validate URL for HTTP mode
	if config.Mode == "http" && config.URL == "" {
		log.Fatal("URL is required for HTTP mode")
	}

	// Health check based on the selected mode
	var err error
	switch strings.ToLower(config.Mode) {
	case "mysql":
		err = mysql.Test(mysql.Config{
			Host:     config.Host,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
			Database: config.Database,
			Timeout:  config.Timeout,
		})
	case "postgres":
		err = postgres.Test(postgres.Config{
			Host:     config.Host,
			Port:     config.Port,
			Username: config.Username,
			Password: config.Password,
			Database: config.Database,
			SSLMode:  config.SSLMode,
			Timeout:  config.Timeout,
		})
	case "redis":
		err = redis.Test(redis.Config{
			Host:     config.Host,
			Port:     config.Port,
			Password: config.Password,
			DB:       config.RedisDB,
			Timeout:  config.Timeout,
		})
	case "port":
		err = port.Test(port.Config{
			Host:    config.Host,
			Port:    config.Port,
			Timeout: config.Timeout,
		})
	case "http":
		err = http.Test(http.Config{
			URL:     config.URL,
			Method:  config.HTTPMethod,
			Timeout: config.Timeout,
		})
	default:
		log.Fatalf("Unsupported mode: %s", config.Mode)
	}

	if err != nil {
		log.Fatalf("Connection test failed: %v", err)
	}

	fmt.Printf("Successfully connected to %s\n", config.Mode)
	os.Exit(0)
} 