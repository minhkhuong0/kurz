APP_NAME := kurz
BUILD_DIR := bin

VERSION ?= $(shell git describe --tags --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildDate=$(DATE)"

.PHONY: api frontend run image fmt clean
server: frontend api

api:
	@echo "Building Go server"
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) .

frontend:
	npm --prefix frontend run build

frontend-install:
	npm --prefix frontend ci

frontend-dev:
	npm --prefix frontend run dev

run:
	@echo "Running Go server locally"
	go run $(LDFLAGS) .

image:
	@echo "Building Docker image"
	docker build . -t kurz:latest

fmt:
	gofmt -w .

clean:
	rm -r bin frontend/dist
