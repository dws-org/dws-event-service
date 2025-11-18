
.PHONY: format
# Run go fmt against code
format:
	go run github.com/steebchen/prisma-client-go format
	go fmt ./...

.PHONY: fmt
# fmt is an alias for format
fmt: format

run:
	go run ./main.go server 

generate:
	go run github.com/steebchen/prisma-client-go generate

dbpush:
	go run github.com/steebchen/prisma-client-go db push

dbpull:
	go run github.com/steebchen/prisma-client-go db pull

