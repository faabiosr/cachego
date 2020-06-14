# Cachego - File driver
The driver stores the cache data in file.

## Usage

```go
package main

import (
	"log"
	"time"

	"github.com/faabiosr/cachego/file"
)

func main() {
	cache := file.New("./cache-dir/")

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
