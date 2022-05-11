package mongo

import (
	"errors"
	"time"

	"github.com/faabiosr/cachego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	mgoCache struct {
		collection *mgo.Collection
	}

	mgoContent struct {
		Duration int64
		Key      string `bson:"_id"`
		Value    string
	}
)

// New creates an instance of Mongo cache driver
// Deprecated: Use NewMongoDriver instead.
func New(collection *mgo.Collection) cachego.Cache {
	return &mgoCache{collection}
}

// Contains checks if cached key exists in Mongo storage
func (m *mgoCache) Contains(key string) bool {
	_, err := m.Fetch(key)
	return err == nil
}

// Delete the cached key from Mongo storage
func (m *mgoCache) Delete(key string) error {
	return m.collection.Remove(bson.M{"_id": key})
}

// Fetch retrieves the cached value from key of the Mongo storage
func (m *mgoCache) Fetch(key string) (string, error) {
	content := &mgoContent{}

	if err := m.collection.Find(bson.M{"_id": key}).One(content); err != nil {
		return "", err
	}

	if content.Duration == 0 {
		return content.Value, nil
	}

	if content.Duration <= time.Now().Unix() {
		_ = m.Delete(key)
		return "", errors.New("cache expired")
	}

	return content.Value, nil
}

// FetchMulti retrieves multiple cached value from keys of the Mongo storage
func (m *mgoCache) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)
	iter := m.collection.Find(bson.M{"_id": bson.M{"$in": keys}}).Iter()
	content := &mgoContent{}

	for iter.Next(content) {
		result[content.Key] = content.Value
	}

	return result
}

// Flush removes all cached keys of the Mongo storage
func (m *mgoCache) Flush() error {
	_, err := m.collection.RemoveAll(bson.M{})
	return err
}

// Save a value in Mongo storage by key
func (m *mgoCache) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &mgoContent{duration, key, value}

	_, err := m.collection.Upsert(bson.M{"_id": key}, content)
	return err
}
