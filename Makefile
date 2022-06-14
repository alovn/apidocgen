.PHONY: test

test:
	go test ./...

lint:
	golangci-lint run ./...

install: lint test
	go install