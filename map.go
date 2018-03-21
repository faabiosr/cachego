package cachego

import (
	"errors"
	"time"
)

type (
	// MapItem structure for managing data and lifetime
	MapItem struct {
		data     string
		duration int64
	}

	// Map store the data in memory without external server
	Map struct {
		storage map[string]*MapItem
	}
)

// NewMap creates an instance of Map cache driver
func NewMap() *Map {
	storage := make(map[string]*MapItem)

	return &Map{storage}
}

func (m *Map) read(key string) (*MapItem, error) {
	item, ok := m.storage[key]

	if !ok {
		return nil, errors.New("Key not found")
	}

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		m.Delete(key)
		return nil, errors.New("Cache expired")
	}

	return item, nil
}

// Contains checks if cached key exists in Map storage
func (m *Map) Contains(key string) bool {
	if _, err := m.Fetch(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from Map storage
func (m *Map) Delete(key string) error {
	delete(m.storage, key)
	return nil
}

// Fetch retrieves the cached value from key of the Map storage
func (m *Map) Fetch(key string) (string, error) {
	item, err := m.read(key)

	if err != nil {
		return "", err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the Map storage
func (m *Map) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := m.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the Map storage
func (m *Map) Flush() error {
	m.storage = make(map[string]*MapItem)
	return nil
}

// Save a value in Map storage by key
func (m *Map) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	item := &MapItem{value, duration}

	m.storage[key] = item

	return nil
}
