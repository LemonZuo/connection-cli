#!/bin/bash

# 设置变量
BINARY="connection-cli"
# 获取最新的 git tag 作为版本号，如果没有则使用默认版本
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "1.0.0")
BUILD_DIR="./build"
CMD_DIR="./cmd"
DOCKER_REPO="zuokaiqi"

# 加载 build.env 文件（如果存在）
if [ -f "build.env" ]; then
  echo "Loading environment variables from build.env"
  source build.env
  
  # 如果环境变量中有定义，则覆盖默认值
  if [ ! -z "$DOCKER_USERNAME" ]; then
    DOCKER_REPO="$DOCKER_USERNAME"
    echo "Using Docker username from build.env: $DOCKER_REPO"
  fi
fi

# 时间测量函数
start_time() {
  START_TIME=$(date +%s)
  echo "⏱️ Starting operation at $(date +'%Y-%m-%d %H:%M:%S')"
}

end_time() {
  END_TIME=$(date +%s)
  DURATION=$((END_TIME - START_TIME))
  MINUTES=$((DURATION / 60))
  SECONDS=$((DURATION % 60))
  echo "⏱️ Operation completed in ${MINUTES}m ${SECONDS}s"
}

# Docker登录函数
docker_login() {
  if [ ! -z "$DOCKER_USERNAME" ] && [ ! -z "$DOCKER_PASSWORD" ]; then
    echo "🔑 Logging in to Docker Hub as $DOCKER_USERNAME"
    echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
    if [ $? -ne 0 ]; then
      echo "❌ Docker login failed. Please check your credentials."
      exit 1
    fi
    echo "✅ Docker login successful"
    DOCKER_LOGGED_IN=true
  else
    echo "⚠️ No Docker credentials found in build.env. Attempting to use existing login..."
    # 检查是否已经登录
    docker info | grep "Username" > /dev/null
    if [ $? -ne 0 ]; then
      echo "❌ Not logged in to Docker Hub. Please provide credentials in build.env or login manually."
      exit 1
    fi
    echo "✅ Using existing Docker login"
    DOCKER_LOGGED_IN=true
  fi
}

# Docker登出函数
docker_logout() {
  if [ "$DOCKER_LOGGED_IN" = true ]; then
    echo "🔒 Logging out from Docker Hub"
    docker logout
    echo "✅ Docker logout successful"
  fi
}

# 显示构建信息
echo "🚀 Building $BINARY version $VERSION"

# 创建构建目录
mkdir -p $BUILD_DIR

# 构建所有平台
build_all() {
  echo "🔨 Building for all platforms..."
  start_time
  
  # Linux (amd64)
  echo "🐧 Building for Linux (amd64)..."
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$VERSION" -o $BUILD_DIR/$BINARY-linux-amd64 $CMD_DIR/
  
  # Linux (arm64)
  echo "🐧 Building for Linux (arm64)..."
  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-X main.Version=$VERSION" -o $BUILD_DIR/$BINARY-linux-arm64 $CMD_DIR/
  
  # macOS (amd64 - Intel)
  echo "🍎 Building for macOS (amd64)..."
  CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$VERSION" -o $BUILD_DIR/$BINARY-darwin-amd64 $CMD_DIR/
  
  # macOS (arm64 - Apple Silicon)
  echo "🍎 Building for macOS (arm64)..."
  CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.Version=$VERSION" -o $BUILD_DIR/$BINARY-darwin-arm64 $CMD_DIR/
  
  # Windows (amd64)
  echo "🪟 Building for Windows (amd64)..."
  CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=$VERSION" -o $BUILD_DIR/$BINARY-windows-amd64.exe $CMD_DIR/

  echo "✅ All builds completed successfully!"
  ls -la $BUILD_DIR
  end_time
}

# 构建当前平台
build_current() {
  echo "🔨 Building for current platform..."
  start_time
  go build -ldflags="-X main.Version=$VERSION" -o $BUILD_DIR/$BINARY $CMD_DIR/
  echo "✅ Build completed!"
  end_time
}

# 构建并推送Docker镜像
build_docker() {
  echo "🐳 Building and pushing Docker image..."
  start_time
  
  # 检查 Docker Buildx 是否可用
  if ! docker buildx version > /dev/null 2>&1; then
    echo "❌ Error: Docker Buildx not available. Please install Docker Buildx."
    exit 1
  fi
  
  # 自动登录Docker
  docker_login
  
  # 创建新的构建器实例（如果不存在）
  if ! docker buildx inspect multi-platform-build > /dev/null 2>&1; then
    echo "🔧 Creating new Docker Buildx builder..."
    docker buildx create --name multi-platform-build --use
  else
    docker buildx use multi-platform-build
  fi
  
  # 构建多平台镜像
  echo "🔨 Building multi-platform Docker image..."
  docker buildx build --platform linux/amd64,linux/arm64 \
    -t $DOCKER_REPO/$BINARY:$VERSION \
    -t $DOCKER_REPO/$BINARY:latest \
    --build-arg VERSION=$VERSION \
    --push \
    --progress=plain \
    .
  
  PUSH_STATUS=$?
  
  # 自动登出Docker
  docker_logout
  
  if [ $PUSH_STATUS -eq 0 ]; then
    echo "✅ Docker image built and pushed successfully: $DOCKER_REPO/$BINARY:$VERSION"
  else
    echo "❌ Failed to build or push Docker image"
    exit 1
  fi
  
  end_time
}

# 构建Docker镜像但不推送
build_docker_local() {
  echo "🐳 Building Docker image locally..."
  start_time
  
  # 自动登录Docker
  docker_login
    
  # 使用环境变量中的用户名
  if [ "$DOCKER_USERNAME" != "$DOCKER_REPO" ]; then
    echo "Using Docker username from build.env: $DOCKER_USERNAME"
    DOCKER_REPO="$DOCKER_USERNAME"
  fi
  
  # 构建镜像
  echo "🔨 Building Docker image..."
  docker build -t $DOCKER_REPO/$BINARY:$VERSION \
    -t $DOCKER_REPO/$BINARY:latest \
    --build-arg VERSION=$VERSION \
    .
  
  BUILD_STATUS=$?
  
  # 如果构建成功，则推送镜像
  if [ $BUILD_STATUS -eq 0 ]; then
    echo "✅ Docker image built successfully: $DOCKER_REPO/$BINARY:$VERSION"
    echo "🚀 Pushing Docker image to repository..."
    
    docker push $DOCKER_REPO/$BINARY:$VERSION
    docker push $DOCKER_REPO/$BINARY:latest
    
    PUSH_STATUS=$?
    if [ $PUSH_STATUS -eq 0 ]; then
      echo "✅ Docker image pushed successfully"
    else
      echo "❌ Failed to push Docker image"
    fi
  else
    echo "❌ Failed to build Docker image"
    docker_logout
    exit 1
  fi
  
  # 自动登出Docker
  docker_logout
  
  end_time
}

# 清理构建目录
clean() {
  echo "🧹 Cleaning build directory..."
  start_time
  rm -rf $BUILD_DIR
  echo "✅ Clean completed!"
  end_time
}

# 显示帮助信息
show_help() {
  echo "📋 Usage: $0 [option]"
  echo "Options:"
  echo "  all          Build for all platforms (Linux, macOS, Windows)"
  echo "  current      Build for current platform only"
  echo "  docker       Build and push multi-platform Docker image"
  echo "  docker-local Build Docker image locally and push to registry"
  echo "  clean        Clean build directory"
  echo "  help         Show this help message"
  echo ""
  echo "🏷️ Current version: $VERSION (from git tag)"
}

# 根据命令行参数执行相应操作
case "$1" in
  "all")
    build_all
    ;;
  "current")
    build_current
    ;;
  "docker")
    build_docker
    ;;
  "docker-local")
    build_docker_local
    ;;
  "clean")
    clean
    ;;
  "help")
    show_help
    ;;
  *)
    # 默认构建所有平台
    if [ -z "$1" ]; then
      build_all
    else
      echo "❌ Unknown option: $1"
      show_help
      exit 1
    fi
    ;;
esac

exit 0 