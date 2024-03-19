PROJECT_NAME := server
GO_SRC_FILES := $(shell find . -type f -name '*.go')

.PHONY: vet
vet:
	go vet .

.PHONY: fmt
fmt:
	gofmt -l -w .

.PHONY: simplify
simplify:
	gofmt -s -l -w .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test -v -race .

.PHONY: bench
bench:
	go test -v -run=B -bench .

.PHONY: build
build:
	go build -v -o $(PROJECT_NAME) .
