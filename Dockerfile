# Build Stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime Stage
FROM alpine:3.20.1

# Install tzdata and ca-certificates
RUN apk update && \
    apk add --no-cache tzdata ca-certificates && \
    rm -rf /var/cache/apk/*

ENV TZ=Asia/Bangkok

# Set working directory in runtime container
WORKDIR /root/

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

EXPOSE 8888

# Run the compiled Go binary
CMD ["./main"]