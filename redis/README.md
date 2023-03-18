# Cachego - Redis driver
The drivers uses [go-redis](https://github.com/go-redis/redis) to store the cache data.

## Usage

```go
package main

import (
	"log"
	"time"

	rd "github.com/go-redis/redis/v8"

	"github.com/faabiosr/cachego/redis"
)

func main() {
	cache := redis.New(
		rd.NewClient(&rd.Options{
			Addr: ":6379",
		}),
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
