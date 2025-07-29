# Multi-stage Dockerfile for go2rtc
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    pkgconfig \
    alsa-lib-dev \
    linux-headers

# Set working directory
WORKDIR /src

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies with proper error handling
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev

# Build the application
RUN CGO_ENABLED=1 \
    GOOS=linux \
    go build \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o go2rtc \
    .

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    ffmpeg \
    python3 \
    py3-pip \
    curl \
    jq \
    ca-certificates \
    tzdata \
    alsa-lib \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1000 go2rtc && \
    adduser -u 1000 -G go2rtc -s /bin/sh -D go2rtc

# Create necessary directories
RUN mkdir -p /config /data && \
    chown -R go2rtc:go2rtc /config /data

# Copy binary from builder stage
COPY --from=builder /src/go2rtc /usr/local/bin/go2rtc
RUN chmod +x /usr/local/bin/go2rtc

# Switch to non-root user
USER go2rtc

# Set working directory
WORKDIR /config

# Expose ports
EXPOSE 1984 8554 8555/tcp 8555/udp

# Environment variables
ENV GO2RTC_CONFIG=/config/go2rtc.yaml

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:1984/api/info || exit 1

# Labels for better metadata
LABEL org.opencontainers.image.title="go2rtc" \
      org.opencontainers.image.description="Ultimate camera streaming application" \
      org.opencontainers.image.source="https://github.com/hamza-farouk/go2rtc" \
      org.opencontainers.image.documentation="https://github.com/hamza-farouk/go2rtc#readme"

# Default command
CMD ["go2rtc"]
