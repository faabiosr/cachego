# Cachego - Chain driver
The chain driver deals with multiple driver at same time, it could save the key in multiple drivers
and for fetching the driver will call the first one, if fails it will try the next until fail.

## Usage

```go
package main

import (
	"log"
	"time"

	bt "go.etcd.io/bbolt"

	"github.com/faabiosr/cachego/bolt"
	"github.com/faabiosr/cachego/chain"
	"github.com/faabiosr/cachego/sync"
)

func main() {
	db, err := bt.Open("cache.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	cache := chain.New(
		bolt.New(db),
		sync.New(),
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
