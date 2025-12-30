# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Install build dependencies (gcc, make, musl-dev for CGO, linux-headers for stdlib)
RUN apk add --no-cache git gcc make musl-dev linux-headers

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application with CGO enabled for SQLite support
RUN PKG_CONFIG_PATH=/usr/lib/pkgconfig \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc \
    go build \
    -ldflags="-w -s" \
    -o /app/bin/api \
    ./cmd/api

# Final stage
FROM alpine:latest

# Install runtime dependencies
# sqlite3 and sqlcipher for database
RUN apk add --no-cache \
    sqlite \
    ca-certificates \
    tzdata \
    && addgroup -g 1000 appuser \
    && adduser -D -u 1000 -G appuser -h /app appuser

# Set timezone
ENV TZ=UTC

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/api /app/api

# Copy configs directory from builder
COPY --from=builder /app/configs /app/configs

# Create directories with proper permissions
RUN mkdir -p /app/data \
    /app/logs \
    && chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/api"]
