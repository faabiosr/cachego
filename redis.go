package cachego

import (
	"gopkg.in/redis.v4"
	"time"
)

type (
	// Redis it's a wrap around the redis driver
	Redis struct {
		driver *redis.Client
	}
)

// NewRedis creates an instance of Redis cache driver
func NewRedis(driver *redis.Client) *Redis {
	return &Redis{driver}
}

// Contains checks if cached key exists in Redis storage
func (r *Redis) Contains(key string) bool {
	status, err := r.driver.Exists(key).Result()

	if err != nil {
		return false
	}

	return status
}

// Delete the cached key from Redis storage
func (r *Redis) Delete(key string) error {
	return r.driver.Del(key).Err()
}

// Fetch retrieves the cached value from key of the Redis storage
func (r *Redis) Fetch(key string) (string, error) {
	value, err := r.driver.Get(key).Result()

	if err != nil {
		return "", err
	}

	return value, nil
}

// FetchMulti retrieves multiple cached value from keys of the Redis storage
func (r *Redis) FetchMulti(keys []string) map[string]string {
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
func (r *Redis) Flush() error {
	return r.driver.FlushAll().Err()
}

// Save a value in Redis storage by key
func (r *Redis) Save(key string, value string, lifeTime time.Duration) error {
	return r.driver.Set(key, value, lifeTime).Err()
}
