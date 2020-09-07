package memcached

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/faabiosr/cachego"
)

type memcached struct {
	driver *memcache.Client
}

// New creates an instance of Memcached cache driver
func New(driver *memcache.Client) cachego.Cache {
	return &memcached{driver}
}

// Contains checks if cached key exists in Memcached storage
func (m *memcached) Contains(key string) bool {
	_, err := m.Fetch(key)
	return err == nil
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

	return string(item.Value[:]), nil
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
	return m.driver.Set(
		&memcache.Item{
			Key:        key,
			Value:      []byte(value),
			Expiration: int32(lifeTime.Seconds()),
		},
	)
}
