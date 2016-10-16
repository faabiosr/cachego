package cachego

import (
	"github.com/bradfitz/gomemcache/memcache"
	"time"
)

type Memcached struct {
	driver *memcache.Client
}

func (m *Memcached) Contains(key string) bool {
	if _, err := m.Fetch(key); err != nil {
		return false
	}

	return true
}

func (m *Memcached) Fetch(key string) (string, error) {
	item, err := m.driver.Get(key)

	if err != nil {
		return "", err
	}

	value := string(item.Value[:])

	return value, nil
}

func (m *Memcached) FetchMulti(keys []string) map[string]string {
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

func (m *Memcached) Save(key string, value string, lifeTime time.Duration) error {
	err := m.driver.Set(
		&memcache.Item{
			Key:        key,
			Value:      []byte(value),
			Expiration: int32(lifeTime.Seconds()),
		},
	)

	return err
}

func (m *Memcached) Delete(key string) error {
	return m.driver.Delete(key)
}

func (m *Memcached) Flush() error {
	return m.driver.FlushAll()
}
