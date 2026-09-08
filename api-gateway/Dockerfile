FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gateway ./cmd/gateway

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/gateway .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./gateway"]