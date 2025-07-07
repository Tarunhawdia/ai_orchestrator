# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 is important for creating statically linked binaries,
# which are easier to run in minimal base images.
# GOOS=linux ensures the binary is built for Linux, even if building on another OS.
# -o main specifies the output binary name.
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./services/orchestrator

# Stage 2: Create a minimal production image
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Set the working directory to /root/ where our binary will be
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]