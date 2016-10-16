package cachego

import (
	"gopkg.in/redis.v4"
	"time"
)

type Redis struct {
	driver *redis.Client
}

func (r *Redis) Fetch(key string) (string, error) {
	value, err := r.driver.Get(key).Result()

	if err != nil {
		return "", err
	}

	return value, nil
}

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

func (r *Redis) Contains(key string) bool {
	status, err := r.driver.Exists(key).Result()

	if err != nil {
		return false
	}

	return status
}

func (r *Redis) Save(key string, value string, lifeTime time.Duration) error {
	return r.driver.Set(key, value, lifeTime).Err()
}

func (r *Redis) Delete(key string) error {
	return r.driver.Del(key).Err()
}

func (r *Redis) Flush() error {
	return r.driver.FlushAll().Err()
}
