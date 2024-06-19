# Build Stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the Go app
COPY . .
RUN go build -o main .

# Final Stage
FROM alpine:3.19
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Remove unnecessary files and dependencies
RUN apk add --no-cache curl && \
    rm -rf /var/cache/apk/*

# Set a non-root user
RUN adduser -D appuser
USER appuser

# Expose port and define health check
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
