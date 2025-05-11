FROM golang:1.23-alpine as builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .

COPY config-prod.toml .

CMD ["./main", "-c", "config-prod.toml"]
