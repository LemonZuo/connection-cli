#!/bin/sh

echo "Connection-CLI Version: $(/app/connection-cli --version)"

# 如果MODE未设置或为help，则运行一个永久循环以保持容器活跃
if [ -z "$MODE" ] || [ "$MODE" = "help" ]; then
    echo "Connection-CLI is running in daemon mode."
    echo "To run tests, set the MODE environment variable to one of: mysql, postgres, redis, port, http"
    echo "Example: docker run -e MODE=port -e HOST=example.com -e PORT=80 connection-cli"

    # 无限循环保持容器运行
    while true; do
        sleep 3600
    done
else
    # 执行实际的连接测试
    exec /app/connection-cli \
        -mode="$MODE" \
        -host="${HOST:-localhost}" \
        -port="${PORT:-0}" \
        -username="$USERNAME" \
        -password="$PASSWORD" \
        -database="$DATABASE" \
        -url="$URL" \
        -timeout="${TIMEOUT:-5s}" \
        -sslmode="${SSLMODE:-disable}" \
        -redis-db="${REDIS_DB:-0}" \
        -http-method="${HTTP_METHOD:-GET}"
fi