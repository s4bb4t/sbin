# Стадия сборки
FROM golang:1.24.4-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY static ./static
RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

# Финальный минимальный образ
FROM alpine:3.18.3

WORKDIR /

COPY --from=builder /app/app /app

EXPOSE 8011

ENTRYPOINT ["/app"]
