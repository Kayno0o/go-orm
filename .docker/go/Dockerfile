FROM golang:1.22-alpine

WORKDIR /app

RUN apk update && apk add --no-cache make
RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["make", "go"]
