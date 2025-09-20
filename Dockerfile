# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 go build -o /svc-template ./cmd/api

# ---------- Runtime stage ----------
FROM alpine:3.20

# Add a non-root user
RUN adduser -D appuser

WORKDIR /home/appuser
COPY --from=builder /svc-template .

USER appuser

EXPOSE 8080
ENTRYPOINT ["./svc-template"]
