.PHONY: build build-all clean docker-build docker-push docker-buildx docker-login docker-logout test

# Binary name
BINARY=connection-cli
# Get version from git tag, fallback to 1.0.0 if no tag exists
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "1.0.0")
DOCKER_REPO=zuokaiqi

# Import environment variables from build.env if it exists
ifneq (,$(wildcard ./build.env))
    include build.env
    # Override Docker repo with username if provided
    ifdef DOCKER_USERNAME
        DOCKER_REPO=$(DOCKER_USERNAME)
    endif
    export
endif

# Build directory
BUILD_DIR=./build

# Time measurement helpers
define start_timer
	@echo "⏱️ Starting operation at $$(date +'%Y-%m-%d %H:%M:%S')"
	@START_TIME=$$(date +%s) && export START_TIME
endef

define end_timer
	@END_TIME=$$(date +%s) && \
	DURATION=$$((END_TIME - START_TIME)) && \
	MINUTES=$$((DURATION / 60)) && \
	SECONDS=$$((DURATION % 60)) && \
	echo "⏱️ Operation completed in $${MINUTES}m $${SECONDS}s"
endef

# Main build
build:
	$(call start_timer)
	@echo "🔨 Building for current platform..."
	go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY) ./cmd/
	@echo "✅ Build completed!"
	$(call end_timer)

# Build for all platforms
build-all:
	$(call start_timer)
	@echo "🔨 Building for all platforms..."
	mkdir -p $(BUILD_DIR)
	@echo "🐧 Building for Linux (amd64)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-linux-amd64 ./cmd/
	@echo "🐧 Building for Linux (arm64)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-linux-arm64 ./cmd/
	@echo "🍎 Building for macOS (amd64)..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 ./cmd/
	@echo "🍎 Building for macOS (arm64)..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 ./cmd/
	@echo "🪟 Building for Windows (amd64)..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe ./cmd/
	@echo "✅ All builds completed successfully!"
	ls -la $(BUILD_DIR)
	$(call end_timer)

# Clean build artifacts
clean:
	$(call start_timer)
	@echo "🧹 Cleaning build directory..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Clean completed!"
	$(call end_timer)

# Docker login helper (uses credentials from build.env if available)
docker-login:
ifdef DOCKER_USERNAME
ifdef DOCKER_PASSWORD
	@echo "🔑 Logging in to Docker Hub as $(DOCKER_USERNAME)"
	@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	@if [ $$? -ne 0 ]; then \
		echo "❌ Docker login failed. Please check your credentials."; \
		exit 1; \
	fi
	@echo "✅ Docker login successful"
else
	@echo "⚠️ DOCKER_PASSWORD not found in build.env. Attempting to use existing login..."
	@docker info | grep "Username" > /dev/null || { echo "❌ Not logged in to Docker Hub. Please provide credentials in build.env or login manually."; exit 1; }
	@echo "✅ Using existing Docker login"
endif
else
	@echo "⚠️ DOCKER_USERNAME not found in build.env. Attempting to use existing login..."
	@docker info | grep "Username" > /dev/null || { echo "❌ Not logged in to Docker Hub. Please provide credentials in build.env or login manually."; exit 1; }
	@echo "✅ Using existing Docker login"
endif

# Docker logout helper
docker-logout:
	@echo "🔒 Logging out from Docker Hub"
	@docker logout
	@echo "✅ Docker logout successful"

# Build Docker image locally
docker-build: 
	$(call start_timer)
	@echo "🐳 Building Docker image locally..."
	@$(MAKE) docker-login
	docker build -t $(DOCKER_REPO)/$(BINARY):$(VERSION) \
		-t $(DOCKER_REPO)/$(BINARY):latest \
		--build-arg VERSION=$(VERSION) \
		.
	@RESULT=$$?; \
	if [ $$RESULT -eq 0 ]; then \
		echo "✅ Docker image built successfully: $(DOCKER_REPO)/$(BINARY):$(VERSION)"; \
		echo "🚀 Pushing Docker image to repository..."; \
		docker push $(DOCKER_REPO)/$(BINARY):$(VERSION); \
		docker push $(DOCKER_REPO)/$(BINARY):latest; \
		PUSH_RESULT=$$?; \
		if [ $$PUSH_RESULT -eq 0 ]; then \
			echo "✅ Docker image pushed successfully"; \
		else \
			echo "❌ Failed to push Docker image"; \
		fi; \
	else \
		echo "❌ Failed to build Docker image"; \
	fi
	@$(MAKE) docker-logout
	$(call end_timer)

# Push Docker image to repository
docker-push: docker-login
	$(call start_timer)
	@echo "🚀 Pushing Docker images to repository..."
	docker push $(DOCKER_REPO)/$(BINARY):$(VERSION)
	docker push $(DOCKER_REPO)/$(BINARY):latest
	@echo "✅ Docker images pushed successfully"
	@$(MAKE) docker-logout
	$(call end_timer)

# Build and push multi-platform Docker images
docker-buildx: 
	$(call start_timer)
	@echo "🐳 Building and pushing multi-platform Docker image..."
	@$(MAKE) docker-login
	docker buildx create --name mybuilder --use || true
	docker buildx use mybuilder
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(DOCKER_REPO)/$(BINARY):$(VERSION) \
		-t $(DOCKER_REPO)/$(BINARY):latest \
		--build-arg VERSION=$(VERSION) \
		--push \
		.
	@RESULT=$$?; \
	if [ $$RESULT -eq 0 ]; then \
		echo "✅ Multi-platform Docker image built and pushed successfully: $(DOCKER_REPO)/$(BINARY):$(VERSION)"; \
	else \
		echo "❌ Failed to build or push multi-platform Docker image"; \
	fi
	@$(MAKE) docker-logout
	$(call end_timer)

# Run the app in Docker
docker-run:
	@echo "🚀 Running container: $(BINARY)"
	docker run -d --name $(BINARY) $(DOCKER_REPO)/$(BINARY):$(VERSION)

# Test running containers
docker-test:
	@echo "🔍 Testing if container is running..."
	@docker ps | grep $(BINARY) || echo "❌ Container not running"

# Test MySQL connection
test-mysql:
	@echo "🔍 Testing MySQL connection..."
	go run -ldflags="-X main.Version=$(VERSION)" ./cmd/main.go -mode=mysql -host=$(HOST) -port=$(PORT) -username=$(USERNAME) -password=$(PASSWORD) -database=$(DATABASE)

# Test PostgreSQL connection
test-postgres:
	@echo "🔍 Testing PostgreSQL connection..."
	go run -ldflags="-X main.Version=$(VERSION)" ./cmd/main.go -mode=postgres -host=$(HOST) -port=$(PORT) -username=$(USERNAME) -password=$(PASSWORD) -database=$(DATABASE) -sslmode=$(SSLMODE)

# Test Redis connection
test-redis:
	@echo "🔍 Testing Redis connection..."
	go run -ldflags="-X main.Version=$(VERSION)" ./cmd/main.go -mode=redis -host=$(HOST) -port=$(PORT) -password=$(PASSWORD) -redis-db=$(REDIS_DB)

# Test port connection
test-port:
	@echo "🔍 Testing port connection..."
	go run -ldflags="-X main.Version=$(VERSION)" ./cmd/main.go -mode=port -host=$(HOST) -port=$(PORT)

# Test HTTP connection
test-http:
	@echo "🔍 Testing HTTP connection..."
	go run -ldflags="-X main.Version=$(VERSION)" ./cmd/main.go -mode=http -url=$(URL) -http-method=$(HTTP_METHOD)

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./... 