package cachego

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
}

func (s *FileTestSuite) SetupTest() {
	directory := "./cache-dir/"

	_ = os.Mkdir(directory, 0777)

	s.cache = NewFile(directory)
	s.assert = assert.New(s.T())
}

func (s *FileTestSuite) TestSaveThrowError() {
	cache := NewFile("./test/")

	s.assert.Regexp("^Unable to save", cache.Save("foo", "bar", 0))
}

func (s *FileTestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 0))
}

func (s *FileTestSuite) TestFetchThrowError() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	cache := NewFile("./test/")

	result, err := cache.Fetch(key)

	s.assert.Regexp("^Unable to open", err)
	s.assert.Empty(result)
}

func (s *FileTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.assert.EqualError(err, "Cache expired")
	s.assert.Empty(result)
}

func (s *FileTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)
	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *FileTestSuite) TestFetchLongCacheDuration() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 10*time.Second)
	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *FileTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *FileTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
	s.assert.Nil(s.cache.Delete("bar"))
}

func (s *FileTestSuite) TestFlushReturnFalseWhenThrowError() {
	cache := NewFile("./test/")

	s.assert.Error(cache.Flush(), "OK")
}

func (s *FileTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *FileTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	cache := NewFile("./test/")
	result := cache.FetchMulti([]string{"foo"})

	s.assert.Len(result, 0)
}

func (s *FileTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *FileTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestFileRunSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}
