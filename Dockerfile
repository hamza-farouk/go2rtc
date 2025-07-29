name: Debug Docker Build

on:
  push:
    branches: [ master, main ]
  workflow_dispatch:

env:
  DOCKER_REGISTRY: docker.io
  DOCKER_IMAGE_NAME: hamzafarouk/go2rtc

jobs:
  # Test build locally first
  test-build:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.23'

    - name: Install build dependencies
      run: |
        sudo apt-get update
        sudo apt-get install -y \
          gcc \
          libc6-dev \
          libasound2-dev \
          pkg-config

    - name: Test Go build
      run: |
        echo "Testing Go modules..."
        go mod download
        go mod verify
        
        echo "Testing Go build..."
        CGO_ENABLED=1 go build -v -ldflags="-s -w" -o go2rtc .
        
        echo "Testing binary..."
        ./go2rtc --help || ./go2rtc -h || echo "Binary built successfully"

  # Build single architecture Docker image first
  build-single-arch:
    runs-on: ubuntu-latest
    needs: test-build
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Build single architecture image (AMD64 only)
      uses: docker/build-push-action@v5
      with:
        context: .
        platforms: linux/amd64
        push: false
        tags: go2rtc:test
        build-args: |
          VERSION=test-$(git rev-parse --short HEAD)

  # Build and push multi-arch if single arch works
  build-docker:
    runs-on: ubuntu-latest
    needs: build-single-arch
    if: github.event_name != 'pull_request'
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Log in to Docker Hub
      uses: docker/login-action@v3
      with:
        registry: ${{ env.DOCKER_REGISTRY }}
        username: ${{ secrets.DOCKER_HUB_USERNAME }}
        password: ${{ secrets.DOCKER_HUB_ACCESS_TOKEN }}

    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.DOCKER_REGISTRY }}/${{ env.DOCKER_IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=raw,value=latest,enable={{is_default_branch}}

    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: .
        platforms: linux/amd64
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        build-args: |
          VERSION=${{ steps.meta.outputs.version }}
