.DEFAULT_GOAL := help

.PHONY: help
help: ## Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
	CGO_ENABLED=0 GOOS=linux go build -o cmd/api/bin/app

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

.PHONY: run
run:
	@echo "Running..."
	@go run cmd/api/main.go

.PHONY: generate-mocks # mockery test mocks defined in .mockery.yaml
generate-mocks:
	mockery