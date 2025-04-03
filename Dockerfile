FROM golang:1.21-alpine AS builder

# Accept build arguments
ARG VERSION=1.0.0

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application for multiple platforms with version information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" -o /app/bin/connection-cli-linux-amd64 ./cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-X main.Version=${VERSION}" -o /app/bin/connection-cli-linux-arm64 ./cmd/
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" -o /app/bin/connection-cli-darwin-amd64 ./cmd/
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.Version=${VERSION}" -o /app/bin/connection-cli-darwin-arm64 ./cmd/
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION}" -o /app/bin/connection-cli-windows-amd64.exe ./cmd/

# Final stage for the main image
FROM alpine:latest

# Accept the version argument in the final stage too
ARG VERSION=1.0.0
ARG TARGETOS
ARG TARGETARCH

# Set version label
LABEL version="${VERSION}" \
      maintainer="zuokaiqi" \
      description="Connection testing CLI tool"

WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/bin/connection-cli--${TARGETOS}-${TARGETARCH} /app/connection-cli

# Make the binary executable
RUN chmod +x /app/connection-cli

# Create a directory to hold all the binaries for distribution
RUN mkdir -p /app/bin
COPY --from=builder /app/bin/* /app/bin/

# Add a script to keep the container running if no mode is specified
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Set the entrypoint to our shell script
ENTRYPOINT ["/app/entrypoint.sh"]
