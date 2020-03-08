package cachego

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// ErrFileOpen returns an error when try to open a file.
const ErrFileOpen = err("unable to open file")

type (
	// File store for caching data
	File struct {
		dir string
	}

	// FileContent it's a structure of cached value
	FileContent struct {
		Duration int64  `json:"duration"`
		Data     string `json:"data,omitempty"`
	}
)

// NewFile creates an instance of File cache
func NewFile(dir string) Cache {
	return &File{dir}
}

func (f *File) createName(key string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))

	filename := hash + ".cachego"

	filePath := filepath.Join(f.dir, filename)

	return filePath
}

func (f *File) read(key string) (*FileContent, error) {
	value, err := ioutil.ReadFile(
		f.createName(key),
	)

	if err != nil {
		return nil, Wrap(ErrFileOpen, err)
	}

	content := &FileContent{}

	err = json.Unmarshal(value, content)

	if err != nil {
		return nil, Wrap(ErrDecode, err)
	}

	if content.Duration == 0 {
		return content, nil
	}

	if content.Duration <= time.Now().Unix() {
		_ = f.Delete(key)
		return nil, ErrCacheExpired
	}

	return content, nil
}

// Contains checks if the cached key exists into the File storage
func (f *File) Contains(key string) bool {

	if _, err := f.read(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from File storage
func (f *File) Delete(key string) error {
	_, err := os.Stat(
		f.createName(key),
	)

	if err != nil && os.IsNotExist(err) {
		return nil
	}

	err = os.Remove(
		f.createName(key),
	)

	if err != nil {
		return Wrap(ErrDelete, err)
	}

	return nil
}

// Fetch retrieves the cached value from key of the File storage
func (f *File) Fetch(key string) (string, error) {
	content, err := f.read(key)

	if err != nil {
		return "", err
	}

	return content.Data, nil
}

// FetchMulti retrieve multiple cached values from keys of the File storage
func (f *File) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := f.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the File storage
func (f *File) Flush() error {
	dir, err := os.Open(f.dir)

	if err != nil {
		return Wrap(ErrFileOpen, err)
	}

	defer dir.Close()

	names, _ := dir.Readdirnames(-1)

	for _, name := range names {
		os.Remove(filepath.Join(f.dir, name))
	}

	return nil
}

// Save a value in File storage by key
func (f *File) Save(key string, value string, lifeTime time.Duration) error {

	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	content := &FileContent{
		duration,
		value,
	}

	data, err := json.Marshal(content)

	if err != nil {
		return Wrap(ErrDecode, err)
	}

	if err := ioutil.WriteFile(f.createName(key), data, 0666); err != nil {
		return Wrap(ErrSave, err)
	}

	return nil
}
