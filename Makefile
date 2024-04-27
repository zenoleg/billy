.PHONY: dc run test lint

dc:
	docker-compose up  --remove-orphans --build

run:
	go run cmd/binomeme/main.go

test:
	go test -race ./...

lint:
	golangci-lint run --timeout=10m
