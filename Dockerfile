FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o bot cmd/bot/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache sqlite

COPY --from=builder /app/bot .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations

CMD goose -dir ./migrations sqlite3 ./data/gymlog.db up && ./bot