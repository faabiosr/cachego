.DEFAULT_GOAL := test

# Clean up
clean:
	@rm -fR ./vendor/ ./cover.*
.PHONY: clean

# Run tests and generates html coverage file
cover: test
	@go tool cover -html=./cover.text -o ./cover.html
	@test -f ./cover.out && rm ./cover.out;
.PHONY: cover

# Up the docker container for testing
docker:
	@docker-compose up -d
.PHONY: docker

# Format all go files
fmt:
	@gofmt -s -w -l $(shell go list -f {{.Dir}} ./...)
.PHONY: fmt

# Run linters
lint:
	@golangci-lint run ./...
.PHONY: lint

# Run tests
test:
	@go test -v -race -coverprofile=./cover.text -covermode=atomic $(shell go list ./...)
.PHONY: test
