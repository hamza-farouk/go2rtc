# Build your custom go2rtc with your edits
FROM golang:1.23-alpine AS builder

# Install basic build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy your modified source code
COPY . .

# Build without CGO to avoid the C compilation issues
# This will exclude ALSA/V4L2 support but include your other modifications
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go2rtc .

# Runtime stage - use the same base as original for compatibility
FROM alpine:3.18

# Install runtime dependencies (same as original)
RUN apk add --no-cache \
    ffmpeg \
    python3 \
    py3-pip \
    curl \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN adduser -D -s /bin/sh go2rtc

# Create config directory
RUN mkdir -p /config && chown go2rtc:go2rtc /config

# Copy YOUR modified binary
COPY --from=builder /app/go2rtc /usr/local/bin/go2rtc
RUN chmod +x /usr/local/bin/go2rtc

# Switch to non-root user  
USER go2rtc

# Set working directory
WORKDIR /config

# Expose ports
EXPOSE 1984 8554 8555

# Environment variables
ENV GO2RTC_CONFIG=/config/go2rtc.yaml

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:1984/api/info || exit 1

# Labels
LABEL org.opencontainers.image.title="go2rtc-custom" \
      org.opencontainers.image.description="go2rtc with custom modifications by hamza-farouk" \
      org.opencontainers.image.source="https://github.com/hamza-farouk/go2rtc"

# Default command
CMD ["go2rtc"]
