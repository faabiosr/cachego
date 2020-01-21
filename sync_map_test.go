package cachego

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SyncMapTestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
}

func (s *SyncMapTestSuite) SetupTest() {
	s.cache = NewSyncMap()
	s.assert = assert.New(s.T())
}

func (s *SyncMapTestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 0))
}

func (s *SyncMapTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.assert.EqualError(err, "Cache expired")
	s.assert.Empty(result)
}

func (s *SyncMapTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *SyncMapTestSuite) TestFetchLongCacheDuration() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 10*time.Second)
	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *SyncMapTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *SyncMapTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
}

func (s *SyncMapTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *SyncMapTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *SyncMapTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestSyncMapRunSuite(t *testing.T) {
	suite.Run(t, new(SyncMapTestSuite))
}
