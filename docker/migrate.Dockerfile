# Build stage
FROM golang:1.12-alpine AS builder
WORKDIR /app
# https://github.com/docker-library/golang/issues/209
RUN apk add --no-cache git
COPY go.mod go.sum ./
COPY tasks/migrate.go tasks/migrate.go
COPY lib lib
COPY migrations migrations
RUN go build -ldflags="-s -w" -o migrate tasks/migrate.go

# TODO set env for migrations dir

# Run stage
FROM alpine AS final
WORKDIR /app
COPY --from=builder /app /app
ENTRYPOINT ./migrate