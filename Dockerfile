FROM golang:1.23.6-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Create a minimal runtime image
FROM alpine:latest
WORKDIR /root/

# Install required dependencies (optional, if needed)
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy the environment file to the final image (if needed)
COPY .env .env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
