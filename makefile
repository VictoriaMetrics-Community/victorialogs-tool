.PHONY: build
# build code
build:
	@echo "Building..."
	@go build -o vtool -v

.PHONY: init
# init env
init:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
# lint code
lint:
	golangci-lint run -v
