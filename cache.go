package cachego

import (
	"time"
)

type Cache interface {
	Contains(key string) bool
	Delete(key string) error
	Fetch(key string) (string, error)
	FetchMulti(keys []string) map[string]string
	Flush() error
	Save(key string, value string, lifeTime time.Duration) error
}
