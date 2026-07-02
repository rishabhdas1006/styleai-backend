# ==========================
# Stage 1 - Build
# ==========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Copy dependency files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o server \
    ./cmd/server

# ==========================
# Stage 2 - Runtime
# ==========================
FROM alpine:3.22

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/server .

# Copy configuration files
COPY configs ./configs

# Cloud Run sets PORT automatically
ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["./server"]