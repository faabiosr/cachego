# Cachego - Mongo driver
The drivers uses [go-mgo](https://github.com/go-mgo/mgo) to store the cache data.

## Usage

```go
package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/faabiosr/cachego/mongo"
)

func main() {
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.Background(), opts)

	cache := mongo.New(
	    client.Database("cache").Collection("cache"),
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
