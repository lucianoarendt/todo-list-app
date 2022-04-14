# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app/todo-list-app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./out/todo-list-app .

EXPOSE 8001

##CMD [ "./out/todo-list-app" ]