# =========================
# 1️⃣ BUILD STAGE
# =========================
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server


# =========================
# 2️⃣ RUNTIME STAGE
# =========================
FROM alpine:3.19

WORKDIR /app

# Install CA cert (important for SMTP, HTTPS, Midtrans)
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/app .

# Create invoices dir
RUN mkdir -p invoices

# Expose port
EXPOSE 8080

# Run app
CMD ["./app"]
