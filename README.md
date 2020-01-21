# Cachego

[![Build Status](https://img.shields.io/travis/faabiosr/cachego/master.svg?style=flat-square)](https://travis-ci.org/faabiosr/cachego)
[![Codecov branch](https://img.shields.io/codecov/c/github/faabiosr/cachego/master.svg?style=flat-square)](https://codecov.io/gh/faabiosr/cachego)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/faabiosr/cachego)
[![Go Report Card](https://goreportcard.com/badge/github.com/faabiosr/cachego?style=flat-square)](https://goreportcard.com/report/github.com/faabiosr/cachego)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/faabiosr/cachego/blob/master/LICENSE)

Simple interface for caching

## Installation

Cachego requires Go 1.9 or later.

```
go get github.com/faabiosr/cachego
```

If you want to get an specific version, please use the example below:

```
go get gopkg.in/faabiosr/cachego.v0
```

## Usage Examples

### Memcached

```go
package main

import (
    "github.com/faabiosr/cachego"
    "github.com/bradfitz/gomemcache/memcache"
)

var cache cachego.Cache

func init() {
    cache = cachego.NewMemcached(memcached.New("localhost:11211"))
}
```

### Redis

```go
package main

import (
    "github.com/faabiosr/cachego"
    "gopkg.in/redis.v4"
)

var cache cachego.Cache

func init() {
    cache = cachego.NewRedis(
        redis.NewClient(&redis.Options{
            Addr: ":6379",
        }),
    )
}
```

### File

```go
package main

import (
    "github.com/faabiosr/cachego"
)

var cache cachego.Cache

func init() {
    cache = cachego.NewFile(
        "/cache-dir/",
    )
}
```

### Map

```go
package main

import (
    "github.com/faabiosr/cachego"
)

var cache cachego.Cache

func init() {
    cache = NewMap()
}
```

### MongoDB

```go
package main

import (
    "github.com/faabiosr/cachego"
    "gopkg.in/mgo.v2"
)

var cache cachego.Cache

func init() {
    session, _ := mgo.Dial(address)

    cache = cachego.NewMongo(
        session.DB("cache").C("cache"),
    )
}
```

### Sqlite3

```go
package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var cache cachego.Cache

func init() {
	db, _ := sql.Open("sqlite3", "./cache.db")

	cache, _ = NewSqlite3(db, "cache")
}
```

### SyncMap

```go
package main

import (
    "github.com/faabiosr/cachego"
)

var cache cachego.Cache

func init() {
    cache = NewSyncMap()
}
```

### BoltDB

```go
package main

import (
    "github.com/faabiosr/cachego"
    bolt "github.com/coreos/bbolt"
)

var cache cachego.Cache

func init() {
    db, _ := bolt.Open("cache.db", 0600, nil)
    cache = NewBolt(db)
}
```

### Chain

```go
package main

import (
    "github.com/faabiosr/cachego"
)

var cache cachego.Cache

func init() {
    memcached := cachego.NewMemcached(
        memcached.New("localhost:11211"),
    )

    redis := cachego.NewRedis(
        redis.NewClient(&redis.Options{
            Addr: ":6379",
        }),
    )

    file := cachego.NewFile(
        "/cache-dir/"
    )

    cache = cachego.NewChain(
        cachego.NewMap(),
        memcached,
        redis,
        file,
    )
}
```

### Usage

```go
package main

import (
    "github.com/faabiosr/cachego"
    "github.com/bradfitz/gomemcache/memcache"
)

func main() {
    cache.Save("foo", "bar")
    cache.Save("john", "doe")

    value, err := cache.Fetch("foo")

    multiple := cache.FetchMulti([]string{"foo", "john"})

    if cache.Contains("foo") {
        cache.Delete("foo")
    }

    cache.Flush()
}
```

## Documentation

Read the full documentation at [https://godoc.org/github.com/faabiosr/cachego](https://godoc.org/github.com/faabiosr/cachego).

## Development

### Requirements

- Install [docker](https://docs.docker.com/install/) and [docker-compose](https://docs.docker.com/compose/install/)
- Install [go dep](https://github.com/golang/dep)

### Makefile
```sh
// Clean up
$ make clean

//Run tests and generates html coverage file
make cover

// Up the docker containers for testing
make docker

// Format all go files
make fmt

//Run linters
make lint

// Run tests
make test
```

## License

This project is released under the MIT licence. See [LICENSE](https://github.com/faabiosr/cachego/blob/master/LICENSE) for more details.
