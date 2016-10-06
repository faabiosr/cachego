package cachego

import (
	"time"
)

type Cache interface {
	Contains(key string) bool
	Delete(key string) bool
	Fetch(key string) (string, bool)
	Flush() bool
	Save(key string, value string, lifeTime time.Duration) bool
}
