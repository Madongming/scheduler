.PHONY: test build-check-race build run

vet:
	go vet ./...

test:
	go clean -testcache && go test ./... -v -cover

build-check-race:
	go build -race -o bin/schedule ./cmd/main.go

build:
	go build -o bin/schedule ./cmd/main.go

run:
	go run ./cmd/main.go
