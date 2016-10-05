package cachego

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcached struct {
	driver *memcache.Client
}

func (m *Memcached) Contains(key string) bool {
	_, status := m.Fetch(key)

	return status
}

func (m *Memcached) Fetch(key string) (string, bool) {
	item, err := m.driver.Get(key)

	if err != nil {
		return "", false
	}

	value := string(item.Value[:])

	return value, true
}

func (m *Memcached) Save(key string, value string) bool {
	err := m.driver.Set(
		&memcache.Item{
			Key:   key,
			Value: []byte(value),
		},
	)

	if err != nil {
		return false
	}

	return true
}

func (m *Memcached) Delete(key string) bool {
	err := m.driver.Delete(key)

	if err != nil {
		return false
	}

	return true
}

func (m *Memcached) Flush() bool {
	err := m.driver.FlushAll()

	if err != nil {
		return false
	}

	return true
}
