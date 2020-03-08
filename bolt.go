package cachego

import (
	"encoding/json"
	"time"

	bolt "github.com/coreos/bbolt"
	errors "github.com/pkg/errors"
)

var (
	boltBucket = []byte("cachego")

	// ErrBoltBucketNotFound returns an error when bucket not found
	ErrBoltBucketNotFound = errors.New("Bucket not found")

	// ErrBoltCacheExpired returns an error when the cache key was expired
	ErrBoltCacheExpired = errors.New("Cache expired")

	// ErrBoltDecodeJSON returns json decoding error message
	ErrBoltDecodeJSON = "Unable to decode json data"

	// ErrBoltFlush returns flush error message
	ErrBoltFlush = "Unable to flush"

	// ErrBoltSave returns save error message
	ErrBoltSave = "Unable to save"
)

type (
	// Bolt store for caching data
	Bolt struct {
		db *bolt.DB
	}

	// BoltContent it's a structure of cached value
	BoltContent struct {
		Duration int64  `json:"duration"`
		Data     string `json:"data,omitempty"`
	}
)

// NewBolt creates an instance of BoltDB cache
func NewBolt(db *bolt.DB) Cache {
	return &Bolt{db}
}

func (b *Bolt) read(key string) (*BoltContent, error) {
	var value []byte

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBucket)

		if bucket == nil {
			return ErrBoltBucketNotFound
		}

		value = bucket.Get([]byte(key))

		return nil
	})

	if err != nil {
		return nil, err
	}

	content := &BoltContent{}

	err = json.Unmarshal(value, content)

	if err != nil {
		return nil, errors.Wrap(err, ErrBoltDecodeJSON)
	}

	if content.Duration == 0 {
		return content, nil
	}

	if content.Duration <= time.Now().Unix() {
		_ = b.Delete(key)
		return nil, ErrBoltCacheExpired
	}

	return content, err
}

// Contains checks if the cached key exists into the BoltDB storage
func (b *Bolt) Contains(key string) bool {
	if _, err := b.read(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from BoltDB storage
func (b *Bolt) Delete(key string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBucket)

		if bucket == nil {
			return ErrBoltBucketNotFound
		}

		return bucket.Delete([]byte(key))
	})

	return err
}

// Fetch retrieves the cached value from key of the BoltDB storage
func (b *Bolt) Fetch(key string) (string, error) {
	content, err := b.read(key)

	if err != nil {
		return "", err
	}

	return content.Data, nil
}

// FetchMulti retrieve multiple cached values from keys of the BoltDB storage
func (b *Bolt) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := b.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the BoltDB storage
func (b *Bolt) Flush() error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(boltBucket)

		if err != nil {
			return errors.Wrap(err, ErrBoltFlush)
		}

		return err
	})

	return err
}

// Save a value in BoltDB storage by key
func (b *Bolt) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &BoltContent{
		duration,
		value,
	}

	data, err := json.Marshal(content)

	if err != nil {
		return errors.Wrap(err, ErrBoltDecodeJSON)
	}

	err = b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(boltBucket)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), data)
	})

	if err != nil {
		return errors.Wrap(err, ErrBoltSave)
	}

	return nil
}
