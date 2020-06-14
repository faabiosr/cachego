# Cachego - SQLite3 driver
The drivers uses [go-sqlite3](https://github.com/mattn/go-sqlite3) to store the cache data.

## Usage

```go
package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/faabiosr/cachego/sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./cache.db")
	if err != nil {
		log.Fatal(err)
	}

	cache, err := sqlite3.New(db, "cache")
	if err != nil {
		log.Fatal(err)
	}

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
