.PHONY: clean configure test test-coverage

.DEFAULT_GOAL := test

clean:
	@rm -fR vendor/ cover.*
	@docker-compose rm -s -f

configure:
	@dep ensure -v
	@docker-compose up -d

test:
	@go test -v .

test-coverage:
	@go test -coverprofile=cover.out -v .
	@go tool cover -html=cover.out -o cover.html
