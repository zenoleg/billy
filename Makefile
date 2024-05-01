.PHONY: dc run test lint

dc:
	docker-compose up  --remove-orphans --build

run:
	go build -o billy cmd/billy/main.go && ./billy

test:
	go test -race ./...

lint:
	golangci-lint run --timeout=10m
