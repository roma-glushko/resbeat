COMMIT ?= $(shell git describe --dirty --long --always)
VERSION := $(shell cat ./VERSION)
LDFLAGS_COMMON := -X main.commitSha=$(COMMIT) -X main.version=$(VERSION) -s -w

run: ##
	@go mod tidy
	@go vet
	@go run main.go

build: ## Build a binary
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o ./dist/resbeat

sandbox-cgroupv1:
	@docker run --cpus="2" --memory="256m" -d alpine:3.18 sleep 1200
