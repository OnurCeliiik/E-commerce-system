FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o product-service main.go

FROM alpine:latest

COPY --from=builder /app/product-service /usr/local/bin/product-service

EXPOSE 8080

CMD ["product-service"]
