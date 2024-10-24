.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: build
build:
	go build -v ./...

.PHONY: all
all: lint test build