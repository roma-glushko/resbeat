GOOS?=darwin
GOARCH?=amd64

CGO_ENABLED?=0
CLI_VERSION_PACKAGE := main
COMMIT ?= $(shell git describe --dirty --long --always --abbrev=15)
VERSION := $(shell cat ./VERSION)
CGO_LDFLAGS_ALLOW := "-Wl,--unresolved-symbols=ignore-in-object-files"
LDFLAGS_COMMON := "-s -w -X $(CLI_VERSION_PACKAGE).commitSha=$(COMMIT) -X $(CLI_VERSION_PACKAGE).version=$(VERSION)"

run: ## Run the app
	@go mod tidy
	@go vet
	@go run main.go

build: ## Build a binary
	@CGO_LDFLAGS_ALLOW=$(CGO_LDFLAGS_ALLOW) CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags $(LDFLAGS_COMMON) -o ./dist/resbeat

build-debug: ## Build with outputting compliler's notes
	@CGO_LDFLAGS_ALLOW=$(CGO_LDFLAGS_ALLOW) CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags $(LDFLAGS_COMMON) -gcflags "-m=2" -o ./dist/resbeat

lint: # Lint the source code
	@echo "🧹 Vetting go.mod.."
	@go vet ./...
	@echo "🧹 Cleaning go.mod.."
	@go mod tidy
	@echo "🧹 GoCI Lint.."
	@golangci-lint run ./...

generate:
	go generate ./...

image:
	@docker build --tag romahlushko/resbeat .

COVERAGE_FILE ?= coverage.out

test: ## Run all tests
	@CGO_ENABLED=$(CGO_ENABLED) go test -v -count=1 -race -shuffle=on -covermode=atomic -coverprofile=$(COVERAGE_FILE) ./...

test-coverage:
	@CGO_ENABLED=$(CGO_ENABLED) go tool cover -html=$(COVERAGE_FILE)

benchmark: ## Run built-in benchmarks
	@go test -v -shuffle=on -run=- -bench=. -benchtime=1x ./...

release-local:  # Perform all artifacts building locally without releasing them actually
	@goreleaser release --snapshot --clean

sandbox: image
	@docker run -p 8000:8000 --cpus="0.5" --memory="150m" --name resbeat-sandbox -d romahlushko/resbeat:latest

linux-%: image-build
	@docker run --rm -v "$(PWD)":/service -w /service -e GOOS=linux romahlushko/resbeat-build:latest make $*
