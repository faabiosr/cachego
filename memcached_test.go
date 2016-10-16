package cachego

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net"
	"testing"
)

type MemcachedTestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
}

func (s *MemcachedTestSuite) SetupTest() {
	address := "localhost:11211"
	conn := memcache.New(address)

	if _, err := net.Dial("tcp", address); err != nil {
		s.T().Skip()
	}

	s.cache = &Memcached{conn}

	s.assert = assert.New(s.T())
}

func (s *MemcachedTestSuite) TestSaveReturnFalseWhenThrowError() {
	memcached := memcache.New("127.0.0.1:22222")

	cache := &Memcached{memcached}

	s.assert.Error(cache.Save("foo", "bar", 0))
}

func (s *MemcachedTestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 0))
}

func (s *MemcachedTestSuite) TestFetchReturnFalseWhenThrowError() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	memcached := memcache.New("127.0.0.1:22222")
	cache := &Memcached{memcached}

	result, err := cache.Fetch(key)

	s.assert.Error(err)
	s.assert.Empty(result)
}

func (s *MemcachedTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *MemcachedTestSuite) TestContains() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *MemcachedTestSuite) TestDeleteReturnFalseWhenThrowError() {
	s.assert.Error(s.cache.Delete("bar"))
}

func (s *MemcachedTestSuite) TestDelete() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
}

func (s *MemcachedTestSuite) TestFlushReturnFalseWhenThrowError() {
	memcached := memcache.New("127.0.0.1:22222")

	cache := &Memcached{memcached}

	s.assert.Error(cache.Flush())
}

func (s *MemcachedTestSuite) TestFlush() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *MemcachedTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	cache := &Memcached{
		memcache.New("127.0.0.1:22222"),
	}

	result := cache.FetchMulti([]string{"foo"})

	s.assert.Len(result, 0)
}

func (s *MemcachedTestSuite) TestFetchMulti() {
	s.cache.Save("foo", "bar", 0)
	s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *MemcachedTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestMemcachedRunSuite(t *testing.T) {
	suite.Run(t, new(MemcachedTestSuite))
}
