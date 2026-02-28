.PHONY: run build dev tidy test

run:
	go run main.go

build:
	go build -o bin/linkbee main.go

dev:
	air

tidy:
	go mod tidy

test:
	go test ./... -v

lint:
	golangci-lint run
