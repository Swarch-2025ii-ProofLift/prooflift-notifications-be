FROM golang:1.25.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o notifications-service ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notifications-service .

CMD ["./notifications-service"]
