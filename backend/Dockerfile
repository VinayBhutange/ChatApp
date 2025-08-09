# Start from the official Go image
FROM golang:1.23-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies with certificate verification disabled
RUN GOPROXY=direct GOSUMDB=off go mod download

# Copy the source code
COPY . .

# Build the application with certificate verification disabled
RUN CGO_ENABLED=1 GOOS=linux GOPROXY=direct GOSUMDB=off go build -a -installsuffix cgo -o chatapp ./cmd/server

# Use a smaller image for the final container
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata sqlite sqlite-dev gcc musl-dev

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/chatapp .

# Expose the application port
EXPOSE 8082

# Command to run the application
CMD ["./chatapp"]
