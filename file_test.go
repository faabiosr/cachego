package cachego

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite

	cache Cache
}

func (s *FileTestSuite) SetupTest() {
	directory := "./cache-dir/"

	_ = os.Mkdir(directory, 0777)

	s.cache = NewFile(directory)
}

func (s *FileTestSuite) TestSaveThrowError() {
	cache := NewFile("./test/")

	s.Assert().Regexp("^Unable to save", cache.Save("foo", "bar", 0))
}

func (s *FileTestSuite) TestSave() {
	s.Assert().Nil(s.cache.Save("foo", "bar", 0))
}

func (s *FileTestSuite) TestFetchThrowError() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	cache := NewFile("./test/")

	result, err := cache.Fetch(key)

	s.Assert().Regexp("^Unable to open", err)
	s.Assert().Empty(result)
}

func (s *FileTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.Assert().EqualError(err, "Cache expired")
	s.Assert().Empty(result)
}

func (s *FileTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)
	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *FileTestSuite) TestFetchLongCacheDuration() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 10*time.Second)
	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *FileTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().True(s.cache.Contains("foo"))
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *FileTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Delete("foo"))
	s.Assert().False(s.cache.Contains("foo"))
	s.Assert().Nil(s.cache.Delete("bar"))
}

func (s *FileTestSuite) TestFlushReturnFalseWhenThrowError() {
	cache := NewFile("./test/")

	s.Assert().Error(cache.Flush(), "OK")
}

func (s *FileTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Flush())
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *FileTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	cache := NewFile("./test/")
	result := cache.FetchMulti([]string{"foo"})

	s.Assert().Len(result, 0)
}

func (s *FileTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.Assert().Len(result, 2)
}

func (s *FileTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.Assert().Len(result, 1)
}

func TestFileRunSuite(t *testing.T) {
	suite.Run(t, new(FileTestSuite))
}
