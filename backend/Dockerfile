# ─── Build stage ─────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Use only the toolchain bundled in this image; never download a newer one
ENV GOTOOLCHAIN=local

# Build deps
RUN apk add --no-cache git ca-certificates tzdata

# Cache dependency downloads as a separate layer
COPY go.mod go.sum ./
RUN go mod download

# Build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" -o /app/server . \
    && echo "--- binary built ---" \
    && ls -lah /app/server

# ─── Runtime stage ────────────────────────────────────────────────────────────
FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl

COPY --from=builder /app/server /app/server
RUN chmod +x /app/server \
    && echo "--- binary in runtime ---" \
    && ls -lah /app/server

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# No subcommand — main() calls cmd.Execute() directly
ENTRYPOINT ["/app/server"]
