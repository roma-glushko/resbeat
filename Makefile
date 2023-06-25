COMMIT ?= $(shell git describe --dirty --long --always)
VERSION := $(shell cat ./VERSION)
LDFLAGS_COMMON := -X main.commitSha=$(COMMIT) -X main.version=$(VERSION) -s -w

run: ##
	@go mod tidy
	@go vet
	@go run main.go

build: ## Build a binary
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o ./dist/resbeat

lint: # Lint the source code
	@go vet ./...
	@go mod tidy

image:
	@docker build --tag romahlushko/resbeat .

test: ## Run all tests
	@go test -v -count=1 -race -shuffle=on -coverprofile=coverage.txt ./...

benchmark: ## Run built-in benchmarks
	@go test -v -shuffle=on -run=- -bench=. -benchtime=1x ./...

release-local:  # Perform all artifacts building locally without releasing them actually
	@goreleaser release --snapshot --clean

sandbox: image
	@docker run -p 8000:8000 --cpus="0.01" --memory="15m" --name resbeat-sandbox -d romahlushko/resbeat:latest
