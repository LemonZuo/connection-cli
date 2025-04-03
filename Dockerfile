FROM golang:1.21-alpine AS builder

# Accept build arguments
ARG VERSION

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application for multiple platforms with version information
RUN go build -ldflags="-s -w -X main.Version=${VERSION} -extldflags '-static'" -o /app/connection-cli ./cmd/

# Final stage for the main image
FROM alpine:latest

WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/connection-cli /app/connection-cli

# Make the binary executable
RUN chmod +x /app/connection-cli

# Add a script to keep the container running if no mode is specified
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Set the entrypoint to our shell script
ENTRYPOINT ["/app/entrypoint.sh"]
