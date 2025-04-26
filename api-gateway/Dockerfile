# api-gateway/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy kode sumber
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./cmd/main.go

# Image tahap kedua untuk runtime yang lebih ringan
FROM alpine:3.18

WORKDIR /app

# Install dependensi yang diperlukan
RUN apk --no-cache add ca-certificates tzdata

# Salin binary dari builder
COPY --from=builder /app/api-gateway .

# Expose port
EXPOSE 8000

# Command untuk menjalankan aplikasi
CMD ["./api-gateway"]