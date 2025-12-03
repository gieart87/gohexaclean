# Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git make protobuf-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build HTTP server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-server ./cmd/http

# Build gRPC server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grpc-server ./cmd/grpc

# Build Worker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# Final stage for HTTP
FROM alpine:latest AS http

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/http-server .
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./http-server"]

# Final stage for gRPC
FROM alpine:latest AS grpc

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/grpc-server .
COPY --from=builder /app/config ./config

EXPOSE 50051

CMD ["./grpc-server"]

# Final stage for Worker
FROM alpine:latest AS worker

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary
COPY --from=builder /app/worker .

CMD ["./worker"]
