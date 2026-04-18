FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o wynnmounts .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/wynnmounts .
RUN mkdir -p /data
ENV WYNNMOUNTS_DB=/data/wynnmounts.db
EXPOSE 8080
ENTRYPOINT ["./wynnmounts", "serve"]
