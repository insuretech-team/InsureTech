# Thin image — binary pre-built via WSL: CGO_ENABLED=0 GOOS=linux go build -o bin/docgen/server
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy pre-built Linux binary
COPY bin/docgen/server ./server

# Copy runtime configs and secrets
COPY ops ./ops
COPY go.mod ./go.mod
COPY backend/inscore/configs ./backend/inscore/configs
COPY backend/inscore/configs ./inscore/configs
COPY backend/inscore/secrets ./backend/inscore/secrets

# Copy docgen templates and generated assets
COPY backend/inscore/microservices/docgen/sidecar/templates ./backend/inscore/microservices/docgen/sidecar/templates
COPY backend/inscore/templates ./backend/inscore/templates
COPY backend/inscore/generated ./backend/inscore/generated

CMD ["./server"]
