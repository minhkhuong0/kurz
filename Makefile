APP_NAME := kurz
BUILD_DIR := bin

VERSION ?= $(shell git describe --tags --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildDate=$(DATE)"

.PHONY: api run image fmt clean
api:
	@echo "Building Go server"
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

image:
	@echo "Building Docker image"
	docker build . -t kurz:latest

fmt:
	gofmt -w .

clean:
	rm -r bin
