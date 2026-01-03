# Simple single-stage Dockerfile for Prisma compatibility
FROM golang:1.23.1

WORKDIR /app

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate docs and Prisma
RUN swag init -g cmd/root.go --output ./docs
RUN go run github.com/steebchen/prisma-client-go generate

# Build
RUN go build -o main .

EXPOSE 6906

CMD ["./main"]
