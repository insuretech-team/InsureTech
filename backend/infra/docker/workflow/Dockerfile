# Build workflow from source so the container always runs the implemented gRPC server.
FROM golang:1.26.1-alpine AS build

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /out/server \
    ./backend/inscore/microservices/workflow/cmd/server

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy freshly built Linux binary
COPY --from=build /out/server ./server

# Copy runtime configs and secrets
COPY ops ./ops
COPY go.mod ./go.mod
COPY backend/inscore/configs ./backend/inscore/configs
COPY backend/inscore/configs ./inscore/configs
COPY backend/inscore/secrets ./backend/inscore/secrets

CMD ["./server"]
