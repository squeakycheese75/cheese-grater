.PHONY: run-debug, lint, test

SHELL = /bin/zsh

run-debug:
	go run cmd/api/main.go

lint:
	golangci-lint version && golangci-lint run --verbose  -E  misspell    

test:
	go test ./... -v
	