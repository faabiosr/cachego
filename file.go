package cachego

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	errors "github.com/pkg/errors"
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

func (f *File) read(key string) (*FileContent, error) {
	value, err := ioutil.ReadFile(
		f.createName(key),
	)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to open")
	}

	content := &FileContent{}

	err = json.Unmarshal(value, content)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to decode json data")
	}

	if content.Duration == 0 {
		return content, nil
	}

	if content.Duration <= time.Now().Unix() {
		f.Delete(key)
		return nil, errors.New("Cache expired")
	}

	return content, nil
}

func (f *File) Contains(key string) bool {

	if _, err := f.read(key); err != nil {
		return false
	}

	return true
}

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
		return errors.Wrap(err, "Unable to delete")
	}

	return nil
}

func (f *File) Fetch(key string) (string, error) {
	content, err := f.read(key)

	if err != nil {
		return "", err
	}

	return content.Data, nil
}

func (f *File) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := f.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

func (f *File) Flush() error {
	dir, err := os.Open(f.dir)

	if err != nil {
		return errors.Wrap(err, "Unable to open file")
	}

	defer dir.Close()

	names, _ := dir.Readdirnames(-1)

	for _, name := range names {
		os.Remove(filepath.Join(f.dir, name))
	}

	return nil
}

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
		return errors.Wrap(err, "Unable to encode json data")
	}

	if err := ioutil.WriteFile(f.createName(key), data, 0666); err != nil {
		return errors.Wrap(err, "Unable to save")
	}

	return nil
}
