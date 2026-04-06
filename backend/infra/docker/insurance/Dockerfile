# Thin image — binary pre-built via WSL: CGO_ENABLED=0 GOOS=linux go build -o bin/insurance/server
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy pre-built Linux binary
COPY bin/insurance/server ./server

# Copy runtime configs
COPY ops ./ops
COPY go.mod ./go.mod
COPY backend/Inscore/configs ./backend/inscore/configs
COPY backend/Inscore/configs ./inscore/configs

CMD ["./server"]
