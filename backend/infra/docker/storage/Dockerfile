# Dockerfile for storage microservice
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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./backend/inscore/cmd/storage/main.go

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
