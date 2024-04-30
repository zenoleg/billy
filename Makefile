.PHONY: run test lint

run:
	go run cmd/billy/main.go && ./billy

test:
	go test -race ./...

lint:
	golangci-lint run --timeout=10m
