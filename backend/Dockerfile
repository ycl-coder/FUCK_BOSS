# Multi-stage build for Go backend
FROM docker.m.daocloud.io/golang:1.24-alpine AS builder

# Set working directory
WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make protobuf protoc

# Install Go protoc plugins
# Add Go bin to PATH for protoc plugins
ENV PATH=$PATH:/root/go/bin
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate gRPC code
RUN chmod +x scripts/generate.sh && ./scripts/generate.sh

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM docker.m.daocloud.io/alpine:latest

# Install ca-certificates and netcat for health check
RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

# Set timezone
ENV TZ=Asia/Shanghai

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/server .

# Copy config example (user should mount their own config)
COPY --from=builder /build/config/config.example.yaml ./config/config.example.yaml

# Expose gRPC port
EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD nc -z localhost 50051 || exit 1

# Run the server
CMD ["./server"]

