# Cachego - Mongo driver
The drivers uses [go-mgo](https://github.com/go-mgo/mgo) to store the cache data.

## Usage

```go
package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/faabiosr/cachego/mongo"
)

func main() {
	session, _ := mgo.Dial("localhost:27017")

	cache := mongo.New(
		session.DB("cache").C("cache"),
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
