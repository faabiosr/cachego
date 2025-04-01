package mongo

import (
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testKeyMongo   = "foo1"
	testValueMongo = "bar"
)

func TestMongo(t *testing.T) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		t.Skip(err)
	}
	collection := client.Database("cache").Collection("cache")

	cache := New(collection)

	if err := cache.Save(testKeyMongo, testValueMongo, 1*time.Nanosecond); err != nil {
		t.Errorf("save fail: expected nil, got %v", err)
	}

	if v, err := cache.Fetch(testKeyMongo); err == nil {
		t.Errorf("fetch fail: expected an error, got %v value %v", err, v)
	}

	_ = cache.Save(testKeyMongo, testValueMongo, 10*time.Second)

	if res, _ := cache.Fetch(testKeyMongo); res != testValueMongo {
		t.Errorf("fetch fail, wrong value : expected %s, got %s", testValueMongo, res)
	}

	_ = cache.Save(testKeyMongo, testValueMongo, 0)

	if !cache.Contains(testKeyMongo) {
		t.Errorf("contains failed: the key %s should be exist", testKeyMongo)
	}

	_ = cache.Save("bar", testValueMongo, 0)

	if values := cache.FetchMulti([]string{testKeyMongo, "bar"}); len(values) != 2 {
		fmt.Println(values)
		t.Errorf("fetch multi failed: expected %d, got %d", 2, len(values))
	}

	if err := cache.Flush(); err != nil {
		t.Errorf("flush failed: expected nil, got %v", err)
	}

	if cache.Contains(testKeyMongo) {
		t.Errorf("contains failed: the key %s should not be exist", testKeyMongo)
	}
}
