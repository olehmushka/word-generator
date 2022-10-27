SHELL=/bin/sh

install:
	go mod download

test:
	go test ./...

test_force:
	go clean -testcache && go test ./...

test_force_v:
	go clean -testcache && go test -v .//...

test_coverage:
	go test -cover .//...

fmt:
	go fmt ./...

tidy:
	go mod tidy

lint:
	golangci-lint run --timeout 60m
