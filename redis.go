package cachego

import (
	"gopkg.in/redis.v4"
	"time"
)

type Redis struct {
	driver *redis.Client
}

func (r *Redis) Fetch(key string) (string, bool) {
	value, err := r.driver.Get(key).Result()

	if err != nil {
		return "", false
	}

	return value, true
}

func (r *Redis) Contains(key string) bool {
	status, err := r.driver.Exists(key).Result()

	if err != nil {
		return false
	}

	return status
}

func (r *Redis) Save(key string, value string, lifeTime time.Duration) bool {
	err := r.driver.Set(key, value, lifeTime).Err()

	if err != nil {
		return false
	}

	return true
}

func (r *Redis) Delete(key string) bool {
	status, err := r.driver.Del(key).Result()

	if err != nil {
		return false
	}

	if status > 0 {
		return true
	}

	return false
}

func (r *Redis) Flush() bool {
	err := r.driver.FlushAll().Err()

	if err != nil {
		return false
	}

	return true
}
