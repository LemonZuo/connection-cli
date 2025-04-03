# Connection CLI

[English](README.md) | [中文](README_zh.md)

![GitHub License](https://img.shields.io/github/license/LemonZuo/connection-cli)
![GitHub commit activity](https://img.shields.io/github/commit-activity/t/LemonZuo/connection-cli)
![GitHub Release](https://img.shields.io/github/v/release/LemonZuo/connection-cli?color=brightgreen)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/LemonZuo/connection-cli/.github%2Fworkflows%2Frelease.yml)
![Docker Pulls](https://img.shields.io/docker/pulls/zuokaiqi/connection-cli)


A Go application that tests connectivity to various services including MySQL, PostgreSQL, Redis, port checks, and HTTP URLs.

## Features

- MySQL connection testing
- PostgreSQL connection testing
- Redis connection testing (with username/password support)
- TCP port connectivity testing
- HTTP URL accessibility testing
- Multi-platform builds (Linux, macOS, Windows)
- Docker support with health checks
- Command-line interface and environment variable configuration
- Automatic logging to app.log file

## Installation

### Using Go

```bash
go get -u github.com/LemonZuo/connection-cli
```

### Using Docker

```bash
docker pull zuokaiqi/connection-cli:latest
```

### Building from Source

```bash
git clone https://github.com/LemonZuo/connection-cli.git
cd connection-cli
make build          # Build for current platform
make build-all      # Build for all platforms
```

## Usage

### Command-line Arguments

```
Usage of connection-cli:
  -mode string
        Testing mode: mysql, postgres, redis, port, http
  -host string
        Host to connect to (default "localhost")
  -port int
        Port to connect to
  -username string
        Username for database connections
  -password string
        Password for database connections
  -database string
        Database name for database connections
  -url string
        URL for HTTP testing
  -timeout duration
        Connection timeout (default 5s)
  -sslmode string
        SSL mode for PostgreSQL (default "disable")
  -redis-db int
        Redis database number
  -http-method string
        HTTP method for HTTP testing (default "GET")
```

### Environment Variables

You can also use environment variables instead of command-line arguments:

- `MODE`: Testing mode (mysql, postgres, redis, port, http)
- `HOST`: Host to connect to
- `PORT`: Port to connect to
- `USERNAME`: Username for database connections
- `PASSWORD`: Password for database connections
- `DATABASE`: Database name for database connections
- `URL`: URL for HTTP testing
- `TIMEOUT`: Connection timeout
- `SSLMODE`: SSL mode for PostgreSQL
- `REDIS_DB`: Redis database number
- `HTTP_METHOD`: HTTP method for HTTP testing

### Examples

#### MySQL Connection Test

```bash
connection-cli -mode=mysql -host=localhost -port=3306 -username=root -password=secret -database=mydb
```

#### PostgreSQL Connection Test

```bash
connection-cli -mode=postgres -host=localhost -port=5432 -username=postgres -password=secret -database=mydb
```

#### Redis Connection Test

```bash
connection-cli -mode=redis -host=localhost -port=6379 -username=redisuser -password=secret
```

#### Port Connection Test

```bash
connection-cli -mode=port -host=localhost -port=8080
```

#### HTTP URL Test

```bash
connection-cli -mode=http -url=https://example.com
```

### Using with Docker

```bash
# Running as a daemon (will keep container running)
docker run -d --name connection-cli zuokaiqi/connection-cli:latest

# MySQL test
docker run --rm -e MODE=mysql -e HOST=mysql-server -e PORT=3306 -e USERNAME=root -e PASSWORD=secret -e DATABASE=mydb zuokaiqi/connection-cli

# PostgreSQL test
docker run --rm -e MODE=postgres -e HOST=postgres-server -e PORT=5432 -e USERNAME=postgres -e PASSWORD=secret -e DATABASE=mydb zuokaiqi/connection-cli

# Redis test
docker run --rm -e MODE=redis -e HOST=redis-server -e PORT=6379 -e USERNAME=redisuser -e PASSWORD=secret zuokaiqi/connection-cli

# Port test
docker run --rm -e MODE=port -e HOST=service-host -e PORT=8080 zuokaiqi/connection-cli

# HTTP test
docker run --rm -e MODE=http -e URL=https://example.com zuokaiqi/connection-cli
```

## Using as a Health Check

You can use this tool as a health check in your Docker Compose or Kubernetes deployments:

```yaml
healthcheck:
  test: ["CMD", "/app/connection-cli", "-mode=mysql", "-host=db", "-port=3306", "-username=root", "-password=secret", "-database=mydb"]
  interval: 30s
  timeout: 5s
  retries: 3
  start_period: 15s
```

## Logging

The tool automatically logs all operations to the `app.log` file in the current directory. This includes:
- Connection attempts
- Success and failure messages
- Error details

This allows you to keep a record of all connectivity tests even when running in automated environments.

## License

MIT 

## Star History
****
[![Star History Chart](https://api.star-history.com/svg?repos=LemonZuo/connection-cli&type=Date)](https://www.star-history.com/#LemonZuo/connection-cli&Date)