# Image URL to use all building/pushing image targets
IMG ?= roasts-api:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.DEFAULT_GOAL := build

.PHONY: all
all: build

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./... -coverprofile cover.out


##@ Build
.PHONY: build
build: fmt vet ## Build the Roasts API binary.
	go build -o bin/roasts-api ./cmd/app/

.PHONY: run
run: fmt vet ## Run the Roasts API from your host.
	(export ENV=local && cd cmd/app && go run .)

.PHONY: docker-build
docker-build: test ## Build docker image with the Roasts API.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the Roasts API.
	docker push ${IMG}
