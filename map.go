package cachego

import (
	"errors"
	"time"
)

type MapItem struct {
	data     string
	duration int64
}

type Map struct {
	storage map[string]*MapItem
}

func NewMapCache() *Map {
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

func (m *Map) Contains(key string) bool {
	if _, err := m.Fetch(key); err != nil {
		return false
	}

	return true
}

func (m *Map) Delete(key string) error {
	delete(m.storage, key)
	return nil
}

func (m *Map) Fetch(key string) (string, error) {
	item, err := m.read(key)

	if err != nil {
		return "", err
	}

	return item.data, nil
}

func (m *Map) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := m.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

func (m *Map) Flush() error {
	m.storage = make(map[string]*MapItem)
	return nil
}

func (m *Map) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	item := &MapItem{value, duration}

	m.storage[key] = item

	return nil
}
