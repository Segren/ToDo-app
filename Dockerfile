FROM golang:1.23-alpine 

WORKDIR /app

COPY ../go.mod ../go.sum ./

COPY ../.env ../.env

RUN go mod download

COPY ../ ./

RUN go build -o main ./cmd/toDo-app

ENV DB_HOST=db \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=todoapp \
    DB_PORT=5432

EXPOSE 8080

CMD ["./main"]