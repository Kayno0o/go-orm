FROM golang:1.22-alpine as golang

WORKDIR /app/server

RUN apk update && apk add --no-cache inotify-tools make bash
RUN go install github.com/Joker/jade/cmd/jade@latest

COPY go.mod go.sum ./
RUN go get -u github.com/Joker/jade

CMD ["make", "pug"]
