# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY ./cmd/tweets ./cmd/tweets
COPY ./internal ./internal
COPY ./pkg ./pkg

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tweets ./cmd/tweets

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/tweets .

# Command to run the executable
CMD ["./tweets"]
