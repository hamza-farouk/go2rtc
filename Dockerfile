# Multi-stage Dockerfile optimized for GitHub Actions
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Build arguments
ARG VERSION=dev
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# Set working directory
WORKDIR /src

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with cross-compilation support
RUN CGO_ENABLED=1 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    GOARM=${TARGETVARIANT#v} \
    go build -ldflags="-s -w -X main.version=$VERSION" -o go2rtc .

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
