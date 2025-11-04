# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage with JuiceFS
FROM alpine:latest

# Install JuiceFS and dependencies
RUN apk add --no-cache \
    ca-certificates \
    fuse \
    curl \
    bash \
    redis

# Install JuiceFS binary
RUN JFS_VERSION=1.1.0 && \
    curl -L https://github.com/juicedata/juicefs/releases/download/v${JFS_VERSION}/juicefs-${JFS_VERSION}-linux-amd64.tar.gz | tar -xz && \
    mv juicefs /usr/local/bin/ && \
    chmod +x /usr/local/bin/juicefs

WORKDIR /app

# Copy Go binary from builder
COPY --from=builder /app/main .

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Create mount point
RUN mkdir -p /jfs

EXPOSE 8090

ENTRYPOINT ["/entrypoint.sh"]