package file

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/faabiosr/cachego"
)

type (
	file struct {
		dir string
	}

	fileContent struct {
		Duration int64  `json:"duration"`
		Data     string `json:"data,omitempty"`
	}
)

// New creates an instance of File cache
func New(dir string) cachego.Cache {
	return &file{dir}
}

func (f *file) createName(key string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))

	filePath := filepath.Join(f.dir, fmt.Sprintf("%s.cachego", hash))

	return filePath
}

func (f *file) read(key string) (*fileContent, error) {
	value, err := ioutil.ReadFile(
		f.createName(key),
	)

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

	if content.Duration <= time.Now().Unix() {
		_ = f.Delete(key)
		return nil, errors.New("cache expired")
	}

	return content, nil
}

// Contains checks if the cached key exists into the File storage
func (f *file) Contains(key string) bool {
	if _, err := f.read(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from File storage
func (f *file) Delete(key string) error {
	_, err := os.Stat(
		f.createName(key),
	)

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

	return content.Data, nil
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
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &fileContent{
		duration,
		value,
	}

	data, err := json.Marshal(content)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.createName(key), data, 0666)
}
