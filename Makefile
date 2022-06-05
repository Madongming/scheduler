.PHONY: test build-check-race build run

build:
	go build -o bin/schedule ./cmd/main.go

vet:
	go vet ./...

test:
	go clean -testcache && go test ./... -v -cover

build-check-race:
	go build -race -o bin/schedule ./cmd/main.go

run:
	go run ./cmd/main.go

build-docker:
	docker build -f Dockerfile . -t test-schedule:v1

run-docker: build-docker
	docker run -it test-schedule:v1
