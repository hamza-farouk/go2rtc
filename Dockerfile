# Dockerfile for go2rtc with automatic package exclusion
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    gcc \
    musl-dev \
    sed

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build arguments
ARG VERSION=dev

# Create build script that excludes problematic packages
RUN cat > build.sh << 'EOF'
#!/bin/sh
set -e

# Initialize go modules if needed
if [ ! -f go.mod ]; then
    go mod init github.com/hamza-farouk/go2rtc
fi

# Check if problematic packages exist and create exclusion script
echo "Checking for CGO-dependent packages..."

# Create a temporary main.go that excludes problematic imports
if [ -f main.go ]; then
    # Check if the main.go imports the problematic packages
    if grep -q "github.com/hamza-farouk/go2rtc/internal/alsa\|github.com/hamza-farouk/go2rtc/internal/v4l2" main.go; then
        echo "Found CGO-dependent packages, creating modified build..."
        
        # Create a modified version without problematic imports
        sed '/github\.com\/hamza-farouk\/go2rtc\/internal\/alsa/d; /github\.com\/hamza-farouk\/go2rtc\/internal\/v4l2/d' main.go > main_nocgo.go
        
        # Try building with the modified main
        echo "Building without CGO-dependent packages..."
        CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$VERSION" -o go2rtc main_nocgo.go
    else
        echo "No problematic CGO packages found, building normally..."
        go mod tidy
        CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$VERSION" -o go2rtc .
    fi
else
    echo "No main.go found, building with current structure..."
    go mod tidy
    CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$VERSION" -o go2rtc .
fi
EOF

# Make script executable and run it
RUN chmod +x build.sh && ./build.sh

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    ffmpeg \
    python3 \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN adduser -D -s /bin/sh go2rtc

# Create directories
RUN mkdir -p /config /data && chown -R go2rtc:go2rtc /config /data

# Copy binary
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

# Default command
CMD ["go2rtc"]
