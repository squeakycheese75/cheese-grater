.PHONY: run-debug, lint, test

SHELL = /bin/zsh

run-debug:
	go run cmd/api/main.go

build-cli:
	go build -o cheese-grater cmd/cli/main.go 

build-executables:
	./build-executables.sh 1

lint:
	golangci-lint version && golangci-lint run --verbose  -E  misspell    

test:
	go test ./... -v
	