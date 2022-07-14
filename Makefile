GOPATH:=$(shell go env GOPATH)
APP?=image

.PHONY: init
## init: initialize the application
init:
	go mod tidy

.PHONY: build
## build: build the application
build:
	go build -o build/${APP} cmd/main.go

.PHONY: run
## run: run the application
run:
	go run -v -race cmd/main.go

.PHONY: format
## format: format files
format:
	@go install golang.org/x/tools/cmd/goimports@latest
	goimports -local example.com/image -w .
	gofmt -s -w .
	go mod tidy

.PHONY: lint
## lint: check everything's okay
lint:
	golangci-lint run ./...
	go mod verify
