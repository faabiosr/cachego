# Cachego - Sync driver
The drivers uses [golang sync](https://golang.org/pkg/sync/#Map) to store the cache data.

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
}
```
