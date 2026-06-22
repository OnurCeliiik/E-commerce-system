FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o gateway main.go

FROM alpine:latest

COPY --from=builder /app/gateway /usr/local/bin/gateway

EXPOSE 8080

CMD ["gateway"]
