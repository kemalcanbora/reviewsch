# Build stage
FROM golang:1.22-alpine AS builder
# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/coupon_service/main.go

# Final stage
FROM alpine:3.18
# Add non root user
RUN adduser -D -g '' schwarz
# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
# Copy binary from builder
COPY --from=builder /src/main .
# Use non root user
USER schwarz
# Set environment variables
ENV GIN_MODE=release

EXPOSE 8080
CMD ["./main"]