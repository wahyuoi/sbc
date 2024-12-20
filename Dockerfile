# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

# Stage 2: Create a minimal image with the built binary
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add ffmpeg
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"] 