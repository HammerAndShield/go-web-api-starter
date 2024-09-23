# Development stage
FROM golang:1.23.1

WORKDIR /app

# Clear Go module cache and install the new Air version
RUN go clean -modcache && \
    go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project
COPY .. .

# Command to run Air for hot reloading
CMD ["air", "-c", ".air.toml"]