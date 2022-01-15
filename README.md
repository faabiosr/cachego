# Cachego

[![Codecov branch](https://img.shields.io/codecov/c/github/faabiosr/cachego/master.svg?style=flat-square)](https://codecov.io/gh/faabiosr/cachego)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/faabiosr/cachego)
[![Go Report Card](https://goreportcard.com/badge/github.com/faabiosr/cachego?style=flat-square)](https://goreportcard.com/report/github.com/faabiosr/cachego)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/faabiosr/cachego/blob/master/LICENSE)

Simple interface for caching

## Installation

Cachego requires Go 1.15 or later.

```
go get github.com/faabiosr/cachego
```

## Usage

```go
package main

import (
	"log"
	"time"

	"github.com/faabiosr/cachego/sync"
)

func main() {
	cache := sync.New()

	if err := cache.Save("user_id", "1", 10*time.Second); err != nil {
		log.Fatal(err)
	}

	id, err := cache.Fetch("user_id")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("user id: %s \n", id)

	keys := cache.FetchMulti([]string{"user_id", "user_name"})

	for k, v := range keys {
		log.Printf("%s: %s\n", k, v)
	}

	if cache.Contains("user_name") {
		cache.Delete("user_name")
	}

	if _, err := cache.Fetch("user_name"); err != nil {
		log.Printf("%v\n", err)
	}

	if err := cache.Flush(); err != nil {
		log.Fatal(err)
	}
}

```

## Supported drivers

- [Bolt](/bolt)
- [Chain](/chain)
- [File](/file)
- [Memcached](/memcached)
- [Mongo](/mongo)
- [Redis](/redis)
- [Sqlite3](/sqlite3)
- [Sync](/sync)


## Documentation

Read the full documentation at [https://pkg.go.dev/github.com/faabiosr/cachego](https://pkg.go.dev/github.com/faabiosr/cachego).

## Development

### Requirements

- Install [docker](https://docs.docker.com/install/)
- Install [docker-compose](https://docs.docker.com/compose/install/)

### Makefile
```sh
// Clean up
$ make clean

//Run tests and generates html coverage file
$ make cover

// Up the docker containers for testing
$ make docker

// Format all go files
$ make fmt

//Run linters
$ make lint

// Run tests
$ make test
```

## License

This project is released under the MIT licence. See [LICENSE](https://github.com/faabiosr/cachego/blob/master/LICENSE) for more details.
