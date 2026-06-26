FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o order-service main.go

FROM alpine:latest

COPY --from=builder /app/order-service /usr/local/bin/order-service

EXPOSE 8080

CMD ["order-service"]
