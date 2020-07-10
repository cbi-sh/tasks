.PHONY: setup
setup: ## Install all dependencies
	go get -u github.com/sirupsen/logrus

.PHONY: clean
clean: ## Remove temporary files
	go clean

.PHONY: build
build: ## Build a version
	go build -v ./...

.PHONY: install
install: ## Install app locally
	go install

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
