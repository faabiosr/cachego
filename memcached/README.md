# Cachego - Memcached driver
The drivers uses [gomemcache](https://github.com/bradfitz/gomemcache) to store the cache data.

## Usage

```go
package main

import (
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/faabiosr/cachego/memcached"
)

func main() {
	cache := memcached.New(
		memcache.New("localhost:11211"),
	)

	if err := cache.Save("user_id", "1", 10*time.Second); err != nil {
		log.Fatal(err)
	}

	id, err := cache.Fetch("user_id")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("user id: %s \n", id)
}
```
