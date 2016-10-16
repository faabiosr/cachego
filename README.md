# Cachego

[![Build Status](https://img.shields.io/travis/fabiorphp/cachego/master.svg?style=flat-square)](https://travis-ci.org/fabiorphp/cachego)
[![Coverage Status](https://img.shields.io/coveralls/fabiorphp/cachego/master.svg?style=flat-square)](https://coveralls.io/github/fabiorphp/cachego?branch=master)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/fabiorphp/cachego)

Simple interface around cache drivers

## Installation
Cachego requires Go 1.5 or later.
```
go get github.com/fabiorphp/cachego
```

If you want to get an specific version, please use the example bellow:
```
go get gopkg.in/fabiorphp/cachego.v0
```

## Usage Examples

### Memcached

```go
package main

import (
    "github.com/fabiorphp/cachego"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
    cache := &cachego.Memcached{
        memcached.New("localhost:11211")
    }

    cache.Save("foo", "bar")

    value := cache.Fetch("foo")
    ...
}
```

### Redis

```go
package main

import (
    "github.com/fabiorphp/cachego"
	"gopkg.in/redis.v4"
)

func main() {
	s.cache = &cachego.Redis{
	    redis.NewClient(&redis.Options{
            Addr: ":6379",
	    }),
    }

    cache.Save("foo", "bar")

    value := cache.Fetch("foo")
    ...
}
```

### File

```go
package main

import (
    "github.com/fabiorphp/cachego"
)

func main() {
	s.cache = &cachego.File{
        "/cache-dir/",
    }

    cache.Save("foo", "bar")

    value := cache.Fetch("foo")
    ...
}
```

## Documentation

Read the full documentation at [https://godoc.org/github.com/fabiorphp/cachego](https://godoc.org/github.com/fabiorphp/cachego).

## License

This project is released under the MIT licence. See [LICENCE](https://github.com/fabiorphp/cachego/blob/master/LICENSE) for more details.
