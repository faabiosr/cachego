// Package redis providers a cache driver that stores the cache in Redis.
package redis

import (
	"context"
	"time"

	rd "github.com/redis/go-redis/v9"

	"github.com/faabiosr/cachego"
)

type redis struct {
	driver rd.Cmdable
}

// New creates an instance of Redis cache driver
func New(driver rd.Cmdable) cachego.Cache {
	return &redis{driver}
}

// Contains checks if cached key exists in Redis storage
func (r *redis) Contains(key string) bool {
	i, _ := r.driver.Exists(context.Background(), key).Result()
	return i > 0
}

// Delete the cached key from Redis storage
func (r *redis) Delete(key string) error {
	return r.driver.Del(context.Background(), key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *redis) Fetch(key string) (string, error) {
	return r.driver.Get(context.Background(), key).Result()
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *redis) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	items, err := r.driver.MGet(context.Background(), keys...).Result()
	if err != nil {
		return result
	}

	for i := 0; i < len(keys); i++ {
		if items[i] != nil {
			result[keys[i]] = items[i].(string)
		}
	}

	return result
}

// Flush removes all cached keys of the Redis storage
func (r *redis) Flush() error {
	return r.driver.FlushAll(context.Background()).Err()
}

// Save a value in Redis storage by key
func (r *redis) Save(key string, value string, lifeTime time.Duration) error {
	return r.driver.Set(context.Background(), key, value, lifeTime).Err()
}
