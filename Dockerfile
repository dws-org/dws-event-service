# Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Prisma client
RUN go run github.com/steebchen/prisma-client-go generate

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# Port freigeben
EXPOSE 6906

# Run the binary
CMD ["./main"]
