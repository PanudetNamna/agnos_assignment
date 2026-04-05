# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Run stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /app/config/secret.env ./config/secret.env

EXPOSE 8080

CMD ["./server"]
