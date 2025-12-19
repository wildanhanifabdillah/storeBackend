# ======================
# Build stage
# ======================
FROM golang:1.25.4-alpine AS builder

# Install build deps (wajib di alpine)
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server

# ======================
# Runtime stage
# ======================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/app /app/app

EXPOSE 8080

# Non-root user (security best practice)
RUN adduser -D appuser
USER appuser

ENTRYPOINT ["/app/app"]
