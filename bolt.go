package cachego

import (
	"encoding/json"
	"time"

	bt "github.com/coreos/bbolt"
)

var boltBucket = []byte("cachego")

// ErrBoltBucketNotFound returns an error when bucket not found
const ErrBoltBucketNotFound = err("Bucket not found")

type (
	bolt struct {
		db *bt.DB
	}

	boltContent struct {
		Duration int64  `json:"duration"`
		Data     string `json:"data,omitempty"`
	}
)

// NewBolt creates an instance of BoltDB cache
func NewBolt(db *bt.DB) Cache {
	return &bolt{db}
}

func (b *bolt) read(key string) (*boltContent, error) {
	var value []byte

	err := b.db.View(func(tx *bt.Tx) error {
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

	content := &boltContent{}

	err = json.Unmarshal(value, content)

	if err != nil {
		return nil, Wrap(ErrDecode, err)
	}

	if content.Duration == 0 {
		return content, nil
	}

	if content.Duration <= time.Now().Unix() {
		_ = b.Delete(key)
		return nil, ErrCacheExpired
	}

	return content, err
}

// Contains checks if the cached key exists into the BoltDB storage
func (b *bolt) Contains(key string) bool {
	if _, err := b.read(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from BoltDB storage
func (b *bolt) Delete(key string) error {
	err := b.db.Update(func(tx *bt.Tx) error {
		bucket := tx.Bucket(boltBucket)

		if bucket == nil {
			return ErrBoltBucketNotFound
		}

		return bucket.Delete([]byte(key))
	})

	return err
}

// Fetch retrieves the cached value from key of the BoltDB storage
func (b *bolt) Fetch(key string) (string, error) {
	content, err := b.read(key)

	if err != nil {
		return "", err
	}

	return content.Data, nil
}

// FetchMulti retrieve multiple cached values from keys of the BoltDB storage
func (b *bolt) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := b.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the BoltDB storage
func (b *bolt) Flush() error {
	err := b.db.Update(func(tx *bt.Tx) error {
		err := tx.DeleteBucket(boltBucket)

		if err != nil {
			return Wrap(ErrFlush, err)
		}

		return err
	})

	return err
}

// Save a value in BoltDB storage by key
func (b *bolt) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &boltContent{
		duration,
		value,
	}

	data, err := json.Marshal(content)

	if err != nil {
		return Wrap(ErrDecode, err)
	}

	err = b.db.Update(func(tx *bt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(boltBucket)

		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), data)
	})

	if err != nil {
		return Wrap(ErrSave, err)
	}

	return nil
}
