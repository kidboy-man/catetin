# Multi-stage Dockerfile for Catetin API
# Builder
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git build-base
WORKDIR /src

# Copy server-side module file to preserve module layout and run dependency download
COPY server-side/go.mod ./server-side/go.mod
WORKDIR /src/server-side
RUN go mod download

# Copy rest of the repository into the build context
WORKDIR /src
COPY . .

# Build a statically linked binary from the server-side module
WORKDIR /src/server-side
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /catetin-api ./cmd/api

# Runtime
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# create non-root user
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

# Copy binary and migration files into image so app can run migrations at startup
COPY --from=builder /catetin-api ./catetin-api
COPY --from=builder /src/server-side/internal/infrastructure/database/postgresql/migrations ./internal/infrastructure/database/postgresql/migrations

RUN chown -R app:app /app
USER app

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["./catetin-api"]
