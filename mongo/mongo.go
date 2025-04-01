// Package mongo providers a cache driver that stores the cache in MongoDB.
package mongo

import (
	"context"
	"time"

	"github.com/faabiosr/cachego"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	mongoCache struct {
		collection *mongo.Collection
	}

	mongoContent struct {
		Duration int64
		Key      string `bson:"_id"`
		Value    string
	}
)

// New creates an instance of Mongo cache driver
func New(collection *mongo.Collection) cachego.Cache {
	return &mongoCache{collection}
}

// NewMongoDriver alias for New.
func NewMongoDriver(collection *mongo.Collection) cachego.Cache {
	return New(collection)
}

func (m *mongoCache) Contains(key string) bool {
	_, err := m.Fetch(key)
	return err == nil
}

// Delete the cached key from Mongo storage
func (m *mongoCache) Delete(key string) error {
	_, err := m.collection.DeleteOne(context.TODO(), bson.M{"_id": bson.M{"$eq": key}})
	return err
}

// Fetch retrieves the cached value from key of the Mongo storage
func (m *mongoCache) Fetch(key string) (string, error) {
	content := &mongoContent{}
	result := m.collection.FindOne(context.TODO(), bson.M{"_id": bson.M{"$eq": key}})
	if result == nil {
		return "", cachego.ErrCacheExpired
	}
	if result.Err() != nil {
		return "", result.Err()
	}

	err := result.Decode(&content)
	if err != nil {
		return "", err
	}
	if content.Duration == 0 {
		return content.Value, nil
	}

	if content.Duration <= time.Now().Unix() {
		_ = m.Delete(key)
		return "", cachego.ErrCacheExpired
	}
	return content.Value, nil
}

func (m *mongoCache) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	cur, err := m.collection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": keys}})
	if err != nil {
		return result
	}
	defer func() {
		_ = cur.Close(context.Background())
	}()

	content := &mongoContent{}

	for cur.Next(context.Background()) {
		err := cur.Decode(content)
		if err != nil {
			continue
		}

		result[content.Key] = content.Value
	}
	return result
}

// Flush removes all cached keys of the Mongo storage
func (m *mongoCache) Flush() error {
	_, err := m.collection.DeleteMany(context.TODO(), bson.M{})
	return err
}

// Save a value in Mongo storage by key
func (m *mongoCache) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &mongoContent{duration, key, value}
	opts := options.Replace().SetUpsert(true)
	_, err := m.collection.ReplaceOne(context.TODO(), bson.M{"_id": bson.M{"$eq": key}}, content, opts)
	return err
}
