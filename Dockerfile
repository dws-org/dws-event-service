FROM golang:1.22


WORKDIR /app
COPY . /app


EXPOSE 6906
RUN go run github.com/steebchen/prisma-client-go generate
RUN go build -o main .

CMD ["./main", "server"]