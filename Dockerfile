FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o wynnbreeder .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/wynnbreeder .
RUN mkdir -p /data
ENV WYNNBREEDER_DB=/data/wynnbreeder.db
EXPOSE 8080
ENTRYPOINT ["./wynnbreeder", "serve"]
