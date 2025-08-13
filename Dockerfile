# Stage 1: Build Go binary
FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api

# Stage 2: Run binary
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3000

CMD ["./main"]
