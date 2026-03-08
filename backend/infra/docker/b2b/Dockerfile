# Dockerfile for b2b microservice
FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Install dependencies (if any native libs are needed)
RUN apk add --no-cache git

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
# Explicitly copy generated packages, configs, and secrets so they are available
COPY gen ./gen

# Build the service
# Build the real b2b microservice directly.
# The wrapper at backend/inscore/cmd/b2b/main.go uses 'go run' which needs
# source code + Go toolchain at runtime — not suitable for the final alpine image.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" \
    -o /app/server \
    ./backend/inscore/microservices/b2b/cmd/server/main.go

# Final minimal image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .
COPY .env .
COPY go.mod .
COPY ops ./ops
COPY backend/inscore/configs ./backend/inscore/configs
COPY backend/inscore/templates ./backend/inscore/templates
COPY backend/inscore/secrets ./backend/inscore/secrets

# Run the server
CMD ["./server"]
