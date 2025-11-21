# main/backend/Dockerfile
FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN go build -o backend ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/backend ./

EXPOSE 8080
CMD ["./backend"]
