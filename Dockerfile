# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Copy source code
COPY . .

# Build the binary
# -ldflags="-w -s" strips debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o susydb ./cmd/susydb

# Stage 2: Run (minimal image)
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/susydb .

# Expose the default port
EXPOSE 7379

# Run
CMD ["./susydb", "--addr", ":7379"]
