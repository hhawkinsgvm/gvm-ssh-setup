# GVM SSH Setup Tool Makefile

.PHONY: help build test clean docker-build docker-run install lint

# Variables
BINARY_NAME=gvm-ssh
VERSION?=latest
DOCKER_TAG=ghcr.io/hhawkinsgvm/gvm-ssh-setup:$(VERSION)

help: ## Show this help message
	@echo "GVM SSH Setup Tool - Build Commands"
	@echo "=================================="
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
build: ## Build the binary
	go build -ldflags="-s -w" -o $(BINARY_NAME) ./main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	docker rmi $(DOCKER_TAG) 2>/dev/null || true

lint: ## Run linters
	go fmt ./...
	go vet ./...

##@ Docker
docker-build: ## Build Docker image
	docker build -t $(DOCKER_TAG) .

docker-run: ## Run Docker container interactively
	docker run --rm -it \
		-u $$(id -u):$$(id -g) \
		-e REAL_HOME=/hosthome \
		-v $$HOME:/hosthome \
		$(DOCKER_TAG) wizard

##@ Installation
install: build ## Install binary to /usr/local/bin
	sudo cp $(BINARY_NAME) /usr/local/bin/

##@ Usage Examples
example-setup: ## Show example setup command
	@echo "Example setup commands:"
	@echo ""
	@echo "# Interactive wizard:"
	@echo "./$(BINARY_NAME) wizard"
	@echo ""
	@echo "# Non-interactive setup:"
	@echo "./$(BINARY_NAME) setup --account gvm --git-alias gitlab-git --upload-key"
	@echo ""
	@echo "# Check configuration:"
	@echo "./$(BINARY_NAME) check --alias gitlab-git"
	@echo ""
	@echo "# Test connectivity:"
	@echo "./$(BINARY_NAME) test --alias gitlab-git --repo Global-Vision-Media/my-project"