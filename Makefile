.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v -covermode=atomic -coverprofile=cover.out -coverpkg ./... ./...

.PHONY: build
build:
	go build -o bin/todolint ./cmd/todolint

.PHONY: build-plugin
build-plugin:
	CGO_ENABLED=1 go build -o bin/todolint.so -buildmode=plugin ./plugin

.PHONY: build-all
build-all: build
