COMMIT ?= $(shell git describe --dirty --long --always)
VERSION := $(shell cat ./VERSION)
LDFLAGS_COMMON := -X main.commitSha=$(COMMIT) -X main.version=$(VERSION) -s -w

run: ##
	@go mod tidy
	@go vet
	@go run main.go

build: ## Build a binary
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o ./dist/resbeat

build-debug: ## Build with outputting compliler's notes
	@GOOS=darwin GOARCH=amd64 go build -gcflags "-m=2" -o ./dist/resbeat

lint: # Lint the source code
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 GoCI Lint.."
	@golangci-lint run ./...

image:
	@docker build --tag romahlushko/resbeat .

test: ## Run all tests
	@go test -v -count=1 -race -shuffle=on -covermode=atomic -coverprofile=coverage.out ./...

test-coverage:
	@go tool cover -html=coverage.out

benchmark: ## Run built-in benchmarks
	@go test -v -shuffle=on -run=- -bench=. -benchtime=1x ./...

release-local:  # Perform all artifacts building locally without releasing them actually
	@goreleaser release --snapshot --clean

sandbox: image
	@docker run -p 8000:8000 --cpus="0.5" --memory="150m" --name resbeat-sandbox -d romahlushko/resbeat:latest
