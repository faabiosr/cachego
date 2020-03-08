package cachego

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type (
	memcached struct {
		driver *memcache.Client
	}
)

// NewMemcached creates an instance of Memcached cache driver
func NewMemcached(driver *memcache.Client) Cache {
	return &memcached{driver}
}

// Contains checks if cached key exists in Memcached storage
func (m *memcached) Contains(key string) bool {
	if _, err := m.Fetch(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from Memcached storage
func (m *memcached) Delete(key string) error {
	return m.driver.Delete(key)
}

// Fetch retrieves the cached value from key of the Memcached storage
func (m *memcached) Fetch(key string) (string, error) {
	item, err := m.driver.Get(key)

	if err != nil {
		return "", err
	}

	value := string(item.Value[:])

	return value, nil
}

// FetchMulti retrieves multiple cached value from keys of the Memcached storage
func (m *memcached) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	items, err := m.driver.GetMulti(keys)

	if err != nil {
		return result
	}

	for _, i := range items {
		result[i.Key] = string(i.Value[:])
	}

	return result
}

// Flush removes all cached keys of the Memcached storage
func (m *memcached) Flush() error {
	return m.driver.FlushAll()
}

// Save a value in Memcached storage by key
func (m *memcached) Save(key string, value string, lifeTime time.Duration) error {
	err := m.driver.Set(
		&memcache.Item{
			Key:        key,
			Value:      []byte(value),
			Expiration: int32(lifeTime.Seconds()),
		},
	)

	return err
}
