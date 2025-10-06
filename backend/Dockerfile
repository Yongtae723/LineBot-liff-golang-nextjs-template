# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum
COPY go_pkg/go.mod go_pkg/go.sum /app/go_pkg/
COPY backend/go.mod backend/go.sum /app/backend/

# Download dependencies
RUN cd /app/go_pkg && go mod download
RUN cd /app/backend && go mod download

# Copy source code
COPY go_pkg /app/go_pkg
COPY backend /app/backend

# Build application
RUN cd /app/backend && CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o main ./cmd/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/backend/main .

# Run application
CMD ["./main"]
