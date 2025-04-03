#!/bin/sh

tail -f -n 200 /app/app.log

/app/connection-cli --version

# 如果MODE未设置或为help，则运行一个永久循环以保持容器活跃
if [ -z "$MODE" ] || [ "$MODE" = "help" ]; then
    echo "Connection CLI is running in daemon mode." >> /app/app.log
    echo "To run tests, set the MODE environment variable to one of: mysql, postgres, redis, port, http" >> /app/app.log
    echo "Example: docker run -e MODE=port -e HOST=example.com -e PORT=80 connection-cli" >> /app/app.log
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