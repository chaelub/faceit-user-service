FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN go mod tidy
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init -g internal/api/api.go

ENTRYPOINT go run main.go