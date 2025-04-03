# 连接测试工具 (Connection CLI)

[English](README.md) | [中文](README_zh.md)

一个Go应用程序，用于测试MySQL、PostgreSQL、Redis的连接状态，以及TCP端口连通性检测和HTTP URL访问测试。

## 功能特点

- MySQL连接测试
- PostgreSQL连接测试
- Redis连接测试
- TCP端口连通性测试
- HTTP URL可访问性测试
- 多平台构建支持（Linux、macOS、Windows）
- Docker支持及健康检查
- 命令行界面和环境变量配置

## 安装

### 使用Go安装

```bash
go get -u github.com/LemonZuo/connection-cli
```

### 使用Docker安装

```bash
docker pull zuokaiqi/connection-cli:latest
```

### 从源代码构建

```bash
git clone https://github.com/LemonZuo/connection-cli.git
cd connection-cli
make build          # 构建当前平台版本
make build-all      # 构建所有平台版本
```

## 使用方法

### 命令行参数

```
connection-cli 使用说明:
  -mode string
        测试模式: mysql, postgres, redis, port, http
  -host string
        连接主机 (默认为 "localhost")
  -port int
        连接端口
  -username string
        数据库连接用户名
  -password string
        数据库连接密码
  -database string
        数据库名称
  -url string
        HTTP测试URL
  -timeout duration
        连接超时时间 (默认为 5s)
  -sslmode string
        PostgreSQL的SSL模式 (默认为 "disable")
  -redis-db int
        Redis数据库编号
  -http-method string
        HTTP测试方法 (默认为 "GET")
```

### 环境变量

你也可以使用环境变量代替命令行参数：

- `MODE`: 测试模式 (mysql, postgres, redis, port, http)
- `HOST`: 连接主机
- `PORT`: 连接端口
- `USERNAME`: 数据库连接用户名
- `PASSWORD`: 数据库连接密码
- `DATABASE`: 数据库名称
- `URL`: HTTP测试URL
- `TIMEOUT`: 连接超时时间
- `SSLMODE`: PostgreSQL的SSL模式
- `REDIS_DB`: Redis数据库编号
- `HTTP_METHOD`: HTTP测试方法

### 使用示例

#### MySQL连接测试

```bash
connection-cli -mode=mysql -host=localhost -port=3306 -username=root -password=secret -database=mydb
```

#### PostgreSQL连接测试

```bash
connection-cli -mode=postgres -host=localhost -port=5432 -username=postgres -password=secret -database=mydb
```

#### Redis连接测试

```bash
connection-cli -mode=redis -host=localhost -port=6379 -password=secret
```

#### 端口连接测试

```bash
connection-cli -mode=port -host=localhost -port=8080
```

#### HTTP URL测试

```bash
connection-cli -mode=http -url=https://example.com
```

### 使用Docker

```bash
# 以守护进程模式运行（保持容器运行）
docker run -d --name connection-cli zuokaiqi/connection-cli:latest

# MySQL测试
docker run --rm -e MODE=mysql -e HOST=mysql-server -e PORT=3306 -e USERNAME=root -e PASSWORD=secret -e DATABASE=mydb zuokaiqi/connection-cli

# PostgreSQL测试
docker run --rm -e MODE=postgres -e HOST=postgres-server -e PORT=5432 -e USERNAME=postgres -e PASSWORD=secret -e DATABASE=mydb zuokaiqi/connection-cli

# Redis测试
docker run --rm -e MODE=redis -e HOST=redis-server -e PORT=6379 -e PASSWORD=secret zuokaiqi/connection-cli

# 端口测试
docker run --rm -e MODE=port -e HOST=service-host -e PORT=8080 zuokaiqi/connection-cli

# HTTP测试
docker run --rm -e MODE=http -e URL=https://example.com zuokaiqi/connection-cli
```

## 作为健康检查使用

你可以在Docker Compose或Kubernetes部署中使用此工具作为健康检查：

```yaml
healthcheck:
  test: ["CMD", "/app/connection-cli", "-mode=mysql", "-host=db", "-port=3306", "-username=root", "-password=secret", "-database=mydb"]
  interval: 30s
  timeout: 5s
  retries: 3
  start_period: 15s
```

## 许可证

MIT 

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=LemonZuo/connection-cli&type=Date)](https://www.star-history.com/#LemonZuo/connection-cli&Date)