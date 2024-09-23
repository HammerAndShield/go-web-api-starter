# Build stage
FROM --platform=$BUILDPLATFORM golang:1.23.1 AS builder

ARG TARGETARCH
ARG TARGETOS

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project
COPY .. .

# Build the specific service
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM --platform=$TARGETPLATFORM alpine:latest

RUN apk --no-cache add ca-certificates

# Install curl
RUN apk --no-cache add curl

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]