# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Set ownership of the application directory and binary
RUN chown -R appuser:appgroup /app

# Set executable permissions
RUN chmod +x /app/server

# Expose the port the application will run on
EXPOSE 8080

# Switch to the non-root user
USER appuser

# Command to run the application
CMD ["./server"]