// Package sync providers a cache driver that uses standard golang sync.Map.
package sync

import (
	"errors"
	"sync"
	"time"

	"github.com/faabiosr/cachego"
)

type (
	syncMapItem struct {
		data     string
		duration int64
	}

	syncMap struct {
		storage *sync.Map
	}
)

// New creates an instance of SyncMap cache driver
func New() cachego.Cache {
	return &syncMap{&sync.Map{}}
}

func (sm *syncMap) read(key string) (*syncMapItem, error) {
	v, ok := sm.storage.Load(key)
	if !ok {
		return nil, errors.New("key not found")
	}

	item := v.(*syncMapItem)

	if item.duration == 0 {
		return item, nil
	}

	if item.duration <= time.Now().Unix() {
		_ = sm.Delete(key)
		return nil, cachego.ErrCacheExpired
	}

	return item, nil
}

// Contains checks if cached key exists in SyncMap storage
func (sm *syncMap) Contains(key string) bool {
	_, err := sm.Fetch(key)
	return err == nil
}

// Delete the cached key from SyncMap storage
func (sm *syncMap) Delete(key string) error {
	sm.storage.Delete(key)
	return nil
}

// Fetch retrieves the cached value from key of the SyncMap storage
func (sm *syncMap) Fetch(key string) (string, error) {
	item, err := sm.read(key)
	if err != nil {
		return "", err
	}

	return item.data, nil
}

// FetchMulti retrieves multiple cached value from keys of the SyncMap storage
func (sm *syncMap) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := sm.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the SyncMap storage
func (sm *syncMap) Flush() error {
	sm.storage = &sync.Map{}
	return nil
}

// Save a value in SyncMap storage by key
func (sm *syncMap) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	sm.storage.Store(key, &syncMapItem{value, duration})
	return nil
}
