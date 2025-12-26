FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# KEIN prisma generate hier (lokal/CI machen)

RUN go build -o /app/app .

EXPOSE 6906
CMD ["/app/app"]
