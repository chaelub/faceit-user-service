FROM golang:latest

WORKDIR /faceit-user-service

COPY ./ /faceit-user-service

RUN go mod tidy
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN swag init -g internal/api/api.go
RUN go test

ENTRYPOINT go run main.go