FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy kode 
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o book-service ./cmd/main.go

# Image tahap kedua untuk runtime yang lebih ringan
FROM alpine:3.18

WORKDIR /app

# Install dependensi yang diperlukan
RUN apk --no-cache add ca-certificates tzdata

# Salin binary dari builder
COPY --from=builder /app/book-service .

# Expose port
EXPOSE 8080
EXPOSE 50051

# Command untuk menjalankan aplikasi
CMD ["./book-service"]