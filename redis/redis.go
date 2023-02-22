// Package redis providers a cache driver that stores the cache in Redis.
package redis

import (
	"time"

	rd "gopkg.in/redis.v4"

	"github.com/faabiosr/cachego"
)

type redis struct {
	driver Cmdable
}

// Cmdable defines the base commands used by redis driver.
type Cmdable interface {
	Exists(key string) *rd.BoolCmd
	Del(keys ...string) *rd.IntCmd
	Get(key string) *rd.StringCmd
	MGet(keys ...string) *rd.SliceCmd
	FlushAll() *rd.StatusCmd
	Set(key string, value interface{}, expiration time.Duration) *rd.StatusCmd
}

// New creates an instance of Redis cache driver
func New(driver Cmdable) cachego.Cache {
	return &redis{driver}
}

// Contains checks if cached key exists in Redis storage
func (r *redis) Contains(key string) bool {
	status, _ := r.driver.Exists(key).Result()
	return status
}

// Delete the cached key from Redis storage
func (r *redis) Delete(key string) error {
	return r.driver.Del(key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *redis) Fetch(key string) (string, error) {
	return r.driver.Get(key).Result()
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *redis) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	items, err := r.driver.MGet(keys...).Result()
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
	return r.driver.FlushAll().Err()
}

// Save a value in Redis storage by key
func (r *redis) Save(key string, value string, lifeTime time.Duration) error {
	return r.driver.Set(key, value, lifeTime).Err()
}
