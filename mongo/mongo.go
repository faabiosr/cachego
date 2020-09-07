package mongo

import (
	"errors"
	"time"

	"github.com/faabiosr/cachego"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	mongo struct {
		collection *mgo.Collection
	}

	mongoContent struct {
		Duration int64
		Key      string `bson:"_id"`
		Value    string
	}
)

// New creates an instance of Mongo cache driver
func New(collection *mgo.Collection) cachego.Cache {
	return &mongo{collection}
}

// Contains checks if cached key exists in Mongo storage
func (m *mongo) Contains(key string) bool {
	if _, err := m.Fetch(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from Mongo storage
func (m *mongo) Delete(key string) error {
	return m.collection.Remove(bson.M{"_id": key})
}

// Fetch retrieves the cached value from key of the Mongo storage
func (m *mongo) Fetch(key string) (string, error) {
	content := &mongoContent{}

	err := m.collection.Find(bson.M{"_id": key}).One(content)
	if err != nil {
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
func (m *mongo) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	iter := m.collection.Find(bson.M{"_id": bson.M{"$in": keys}}).Iter()

	content := &mongoContent{}

	for iter.Next(content) {
		result[content.Key] = content.Value
	}

	return result
}

// Flush removes all cached keys of the Mongo storage
func (m *mongo) Flush() error {
	_, err := m.collection.RemoveAll(bson.M{})

	return err
}

// Save a value in Mongo storage by key
func (m *mongo) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &mongoContent{
		duration,
		key,
		value,
	}

	if _, err := m.collection.Upsert(bson.M{"_id": key}, content); err != nil {
		return err
	}

	return nil
}
