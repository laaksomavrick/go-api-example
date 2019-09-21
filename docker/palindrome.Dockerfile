# Build stage
FROM golang:1.12-alpine AS builder
WORKDIR /app
# https://github.com/docker-library/golang/issues/209
RUN apk add --no-cache git
COPY go.mod go.sum ./
COPY src src
COPY lib lib
RUN go build -ldflags="-s -w" -o palindrome src/main.go

# Run stage
FROM alpine AS final
WORKDIR /app
COPY --from=builder /app /app
ENTRYPOINT ./palindrome