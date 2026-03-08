FROM golang:1.25-alpine AS builder

WORKDIR /workspace

RUN apk add --no-cache git

# Cache modules first
COPY go.mod go.sum ./
RUN go mod download

# Copy full source
COPY . .

# Build the dbops binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" \
    -o /app/dbops \
    ./backend/inscore/cmd/dbops/main.go

# ── Final minimal image ──────────────────────────────────────
FROM alpine:latest

WORKDIR /app

# Binary
COPY --from=builder /app/dbops .

# go.mod is required by findProjectRoot() which walks up from os.Getwd()
# looking for go.mod to determine the project root. Without it the binary
# cannot locate migrations/seeds/configs/secrets directories.
COPY go.mod ./

# Runtime data directories (paths resolved relative to projectRoot = /app)
COPY backend/inscore/db/migrations ./backend/inscore/db/migrations
COPY backend/inscore/db/seeds      ./backend/inscore/db/seeds
COPY backend/inscore/configs       ./backend/inscore/configs
COPY backend/inscore/secrets       ./backend/inscore/secrets

CMD ["./dbops", "migrate", "--target=both"]
