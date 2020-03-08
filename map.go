package cachego

import (
	"time"
)

//ErrMapNotFound returns an error when the key is not found.
const ErrMapKeyNotFound = err("key not found")

type (
	mapItem struct {
		data     string
		duration int64
	}

	mapcache struct {
		storage map[string]*mapItem
	}
)

// NewMap creates an instance of Map cache driver
func NewMap() Cache {
	storage := make(map[string]*mapItem)

	return &mapcache{storage}
}

func (m *mapcache) read(key string) (*mapItem, error) {
	item, ok := m.storage[key]

	if !ok {
		return nil, ErrMapKeyNotFound
	}

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		_ = m.Delete(key)
		return nil, ErrCacheExpired
	}

	return item, nil
}

// Contains checks if cached key exists in Map storage
func (m *mapcache) Contains(key string) bool {
	if _, err := m.Fetch(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from Map storage
func (m *mapcache) Delete(key string) error {
	delete(m.storage, key)
	return nil
}

// Fetch retrieves the cached value from key of the Map storage
func (m *mapcache) Fetch(key string) (string, error) {
	item, err := m.read(key)

	if err != nil {
		return "", err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the Map storage
func (m *mapcache) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := m.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the Map storage
func (m *mapcache) Flush() error {
	m.storage = make(map[string]*mapItem)
	return nil
}

// Save a value in Map storage by key
func (m *mapcache) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	item := &mapItem{value, duration}

	m.storage[key] = item

	return nil
}
