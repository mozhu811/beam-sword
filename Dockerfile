# Stage 1: Build the Go binary
FROM golang:1.22.4-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Copy .env file
COPY .env .env

# Build the Go app
RUN go build -o main .

# Stage 2: Run the Go binary
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file and .env from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Install necessary packages if your application requires it, e.g., ca-certificates
RUN apk --no-cache add ca-certificates

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
