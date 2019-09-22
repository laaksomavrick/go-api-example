FROM golang:1.12-alpine
WORKDIR /app
# https://github.com/docker-library/golang/issues/209
RUN apk add --no-cache git
COPY go.mod go.sum ./
COPY lib lib
COPY tests tests

CMD CGO_ENABLED=0 go test -v -count=1 -p=1 ./tests/messages_test.go