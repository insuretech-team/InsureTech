# Dockerfile for media microservice
FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Install dependencies (if any native libs are needed)
RUN apk add --no-cache git

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./backend/inscore/cmd/media/main.go

# Final minimal image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .

# Run the server
CMD ["./server"]
