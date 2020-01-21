.DEFAULT_GOAL := test

# Clean up
clean:
	rm -fR ./cover.*
.PHONY: clean

# Run tests and generates html coverage file
cover: test
	go tool cover -html=./cover.out -o ./cover.html
.PHONY: cover

# Up the docker container for testing
docker:
	docker-compose up -d
.PHONY: docker

# Format all go files
fmt:
	gofmt -s -w -l $(shell go list -f {{.Dir}} ./... | grep -v /vendor/)
.PHONY: fmt

# Run linters
lint:
	golangci-linter run ./...
.PHONY: lint

# Run tests
test:
	go test -v -race -coverprofile=./cover.out $(shell go list ./... | grep -v /vendor/)
.PHONY: test
