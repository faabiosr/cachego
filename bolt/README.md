# Cachego - BoltDB driver
The drivers uses [etcd-io/bbolt](https://github.com/etcd-io/bbolt) to store the cache data.

## Usage

```go
package main

import (
	"log"
	"time"

	bt "go.etcd.io/bbolt"

	"github.com/faabiosr/cachego/bolt"
)

func main() {
	db, err := bt.Open("cache.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	cache := bolt.New(db)

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
