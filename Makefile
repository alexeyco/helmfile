BINARY_NAME := image
BUILD_DIR   := build
CMD         := ./cmd/image

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

.PHONY: mod
mod: ## Tidy and download dependencies
	go mod tidy
	go mod download

.PHONY: fmt
fmt: ## Format code (golangci-lint)
	golangci-lint fmt

.PHONY: lint
lint: ## Run linters (golangci-lint)
	golangci-lint run

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: build
build: ## Build binary into ./build
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD)

.PHONY: run
run: ## Run the app
	go run $(CMD)

.PHONY: generate
generate: ## Generate Dockerfile with latest tool versions
	go run $(CMD)

.PHONY: image
image: generate ## Generate the Dockerfile and build the image locally
	docker build -t ghcr.io/alexeyco/helmfile:local .

.PHONY: mock
mock: ## Generate mocks (uber/mock)
	go generate ./...
