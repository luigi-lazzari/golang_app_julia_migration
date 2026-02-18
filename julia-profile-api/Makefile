.PHONY: help build run test clean swagger docker-build docker-run

SWAG := $(shell which swag 2> /dev/null || echo $(shell go env GOPATH)/bin/swag)


help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/bff-julia-profile-api cmd/api/main.go

run: ## Run the application
	go run cmd/api/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

swagger: ## Generate Swagger documentation
	$(SWAG) init -g cmd/api/main.go -o ./docs

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

docker-build: ## Build Docker image
	docker build -t bff-julia-profile-api:latest .

docker-run: ## Run with docker-compose
	docker-compose up -d

docker-stop: ## Stop docker-compose services
	docker-compose down

docker-logs: ## View docker-compose logs
	docker-compose logs -f

mod-download: ## Download dependencies
	go mod download

mod-tidy: ## Tidy dependencies
	go mod tidy

.DEFAULT_GOAL := help
