# Etapa 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main cmd/web/main.go

# Etapa 2: Runtime
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/main .
COPY entrypoint.sh /app/entrypoint.sh
COPY .env /app/.env 
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080
CMD ["/app/entrypoint.sh", "./main"]
