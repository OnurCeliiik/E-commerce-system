FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o notification-service main.go

FROM alpine:latest

COPY --from=builder /app/notification-service /usr/local/bin/notification-service

EXPOSE 8080

CMD ["notification-service"]
