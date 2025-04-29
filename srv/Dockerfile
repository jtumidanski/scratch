# Build stage
FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/document-storage-service

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/document-storage-service .

# Expose the port
EXPOSE 8080

# Set environment variables
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=document_storage
ENV DB_SSLMODE=disable
ENV PORT=8080

# Run the application
CMD ["/app/document-storage-service"]