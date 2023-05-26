run: ##
	@go mod tidy
	@go vet
	@go run main.go

build: ## Build a binary
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./dist/resbeat
