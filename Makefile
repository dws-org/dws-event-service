
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

# Docker variables - setze diese Variablen oder überschreibe sie beim Aufruf
DOCKER_REGISTRY ?= docker.io
DOCKER_USERNAME ?= your-username
IMAGE_NAME ?= dws-event-service
IMAGE_TAG ?= latest
FULL_IMAGE_NAME = $(DOCKER_REGISTRY)/$(DOCKER_USERNAME)/$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: docker-build
# Build Docker image
docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(FULL_IMAGE_NAME)

.PHONY: docker-buildx
# Build Docker image for multiple architectures (amd64, arm64)
docker-buildx:
	docker buildx create --use --name multiarch-builder 2>/dev/null || true
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(FULL_IMAGE_NAME) \
		--push \
		.

.PHONY: docker-push
# Push Docker image to registry
docker-push: docker-build
	docker push $(FULL_IMAGE_NAME)

.PHONY: docker-buildx-push
# Build and push Docker image for multiple architectures
docker-buildx-push: docker-buildx
	@echo "Image built and pushed for multiple architectures"

.PHONY: docker-login
# Login to Docker registry (für Docker Hub: docker login, für andere: docker login <registry-url>)
docker-login:
	@echo "Bitte melden Sie sich bei Ihrer Registry an:"
	@echo "Für Docker Hub: docker login"
	@echo "Für GitHub Container Registry: docker login ghcr.io"
	@echo "Für andere Registries: docker login <registry-url>"
	@docker login $(DOCKER_REGISTRY)

