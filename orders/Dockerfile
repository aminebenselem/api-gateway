FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod  ./
RUN go mod download

COPY . .

RUN go build -o orders .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/orders .

EXPOSE 8080

CMD ["./orders"]