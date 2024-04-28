.PHONY: run test lint

run:
	go run cmd/binomeme/main.go && ./binomeme

test:
	go test -race ./...

lint:
	golangci-lint run --timeout=10m
