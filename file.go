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

type File struct {
	dir string
}

type FileContent struct {
	Duration int64  `json:"duration"`
	Data     string `json:"data, omitempty"`
}

func (f *File) createName(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	hash := hex.EncodeToString(h.Sum(nil))

	filename := hash + ".cachego"

	filePath := filepath.Join(f.dir, filename)

	return filePath
}

func (f *File) read(key string) (*FileContent, bool) {
	value, err := ioutil.ReadFile(
		f.createName(key),
	)

	if err != nil {
		return nil, false
	}

	content := &FileContent{}

	err = json.Unmarshal(value, content)

	if err != nil {
		return nil, false
	}

	if content.Duration == 0 {
		return content, true
	}

	if content.Duration <= time.Now().Unix() {
		f.Delete(key)
		return nil, false
	}

	return content, true
}

func (f *File) Contains(key string) bool {

	_, ok := f.read(key)

	return ok
}

func (f *File) Delete(key string) bool {
	err := os.Remove(
		f.createName(key),
	)

	if err != nil {
		return false
	}

	return true
}

func (f *File) Fetch(key string) (string, bool) {
	if content, ok := f.read(key); ok {
		return content.Data, true
	}

	return "", false
}

func (f *File) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, ok := f.Fetch(key); ok {
			result[key] = value
		}
	}

	return result
}

func (f *File) Flush() bool {
	dir, err := os.Open(f.dir)

	if err != nil {
		return false
	}

	defer dir.Close()

	names, _ := dir.Readdirnames(-1)

	for _, name := range names {
		os.Remove(filepath.Join(f.dir, name))
	}

	return true
}

func (f *File) Save(key string, value string, lifeTime time.Duration) bool {

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
		return false
	}

	if err := ioutil.WriteFile(f.createName(key), data, 0666); err != nil {
		return false
	}

	return true
}
