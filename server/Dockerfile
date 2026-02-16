# Build stage
FROM golang:1.24.3-alpine AS builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

# Download dependencies (cached as separate layer)
RUN go mod download

COPY . .

# Build with optimizations for production
# -ldflags="-s -w" removes debug info and symbol table to reduce binary size
# CGO_ENABLED=0 for static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -trimpath \
    -o bin/server ./cmd

# Final stage
FROM alpine:3.18
WORKDIR /app

# Add ca-certificates for HTTPS requests and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Use non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appgroup /app/bin/server .

USER appuser

EXPOSE 8080
CMD ["./server"]
