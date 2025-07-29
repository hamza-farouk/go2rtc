# Simple Dockerfile for go2rtc
FROM golang:1.23-alpine AS builder

# Install minimal build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy everything
COPY . .

# Initialize go modules if needed and build
RUN if [ ! -f go.mod ]; then go mod init github.com/hamza-farouk/go2rtc; fi && \
    go mod tidy && \
    CGO_ENABLED=1 go build -o go2rtc

# Runtime stage
FROM alpine:3.18

# Install minimal runtime dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh go2rtc

# Copy binary
COPY --from=builder /app/go2rtc /usr/local/bin/go2rtc
RUN chmod +x /usr/local/bin/go2rtc

# Switch to non-root user
USER go2rtc

# Expose ports
EXPOSE 1984 8554 8555

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1984/ || exit 1

# Default command
CMD ["go2rtc"]
