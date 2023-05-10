// Package file providers a cache driver that stores the cache content in files.
package file

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faabiosr/cachego"
)

type (
	file struct {
		dir string
		sync.RWMutex
	}

	fileContent struct {
		Duration int64  `json:"duration"`
		Data     string `json:"data,omitempty"`
	}
)

const perm = 0o666

// New creates an instance of File cache
func New(dir string) cachego.Cache {
	return &file{dir: dir}
}

func (f *file) createName(key string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))

	return filepath.Join(f.dir, fmt.Sprintf("%s.cachego", hash))
}

func (f *file) read(key string) (*fileContent, error) {
	f.RLock()
	defer f.RUnlock()

	value, err := ioutil.ReadFile(f.createName(key))
	if err != nil {
		return nil, err
	}

	content := &fileContent{}
	if err := json.Unmarshal(value, content); err != nil {
		return nil, err
	}

	if content.Duration == 0 {
		return content, nil
	}

	return content, nil
}

// Contains checks if the cached key exists into the File storage
func (f *file) Contains(key string) bool {
	content, err := f.read(key)
	if err != nil {
		return false
	}

	if f.isExpired(content) {
		_ = f.Delete(key)
		return false
	}
	return true
}

// Delete the cached key from File storage
func (f *file) Delete(key string) error {
	f.Lock()
	defer f.Unlock()

	_, err := os.Stat(f.createName(key))
	if err != nil && os.IsNotExist(err) {
		return nil
	}

	return os.Remove(f.createName(key))
}

// Fetch retrieves the cached value from key of the File storage
func (f *file) Fetch(key string) (string, error) {
	content, err := f.read(key)
	if err != nil {
		return "", err
	}

	if f.isExpired(content) {
		_ = f.Delete(key)
		return "", cachego.ErrCacheExpired
	}

	return content.Data, nil
}

func (f *file) isExpired(content *fileContent) bool {
	return content.Duration > 0 && content.Duration <= time.Now().Unix()
}

// FetchMulti retrieve multiple cached values from keys of the File storage
func (f *file) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := f.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the File storage
func (f *file) Flush() error {
	f.Lock()
	defer f.Unlock()

	dir, err := os.Open(f.dir)
	if err != nil {
		return err
	}

	defer func() {
		_ = dir.Close()
	}()

	names, _ := dir.Readdirnames(-1)

	for _, name := range names {
		_ = os.Remove(filepath.Join(f.dir, name))
	}

	return nil
}

// Save a value in File storage by key
func (f *file) Save(key string, value string, lifeTime time.Duration) error {
	f.Lock()
	defer f.Unlock()

	duration := int64(0)
	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &fileContent{duration, value}
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.createName(key), data, perm)
}
