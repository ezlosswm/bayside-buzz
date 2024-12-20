# Use the official Golang image for building the application
FROM golang:1.23.1 AS builder

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . ./

# Build the application
RUN go build -o main ./

# Use a minimal base image for the runtime
FROM gcr.io/distroless/base-debian11

# Set the working directory in the container
WORKDIR /app

# Copy the compiled application binary from the builder stage
COPY --from=builder /app/main ./

# Expose the application's port (Railway uses PORT environment variable)
EXPOSE 8080

# Use the PORT environment variable provided by Railway
ENV PORT=8080

# Set the entrypoint command
CMD ["./main"]

