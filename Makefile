COMMIT ?= $(shell git describe --dirty --long --always)
VERSION := $(shell cat ./VERSION)
LDFLAGS_COMMON := -X main.commitSha=$(COMMIT) -X main.version=$(VERSION) -s -w

run: ##
	@go mod tidy
	@go vet
	@go run main.go

build: ## Build a binary
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o ./dist/resbeat

image:
	@docker build --tag romahlushko/resbeat .

release-local:  # Perform all artifacts building locally without releasing them actually
	@goreleaser release --snapshot --clean

sandbox: image
	@docker run -p 8000:8000 --cpus="0.01" --memory="15m" --name resbeat-sandbox -d romahlushko/resbeat:latest
