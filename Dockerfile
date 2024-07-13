# Dockerfile

FROM golang:1.22.1-alpine as builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o dns_server ./cmd/dns-server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/dns_server .

COPY config/config.yaml ./config/config.yaml


EXPOSE 53/udp
EXPOSE 443/tcp

CMD ["./dns_server"]