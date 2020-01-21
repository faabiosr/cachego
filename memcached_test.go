package cachego

import (
	"net"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/suite"
)

type MemcachedTestSuite struct {
	suite.Suite
	cache Cache
}

func (s *MemcachedTestSuite) SetupTest() {
	address := "localhost:11211"
	conn := memcache.New(address)

	if _, err := net.Dial("tcp", address); err != nil {
		s.T().Skip()
	}

	s.cache = NewMemcached(conn)
}

func (s *MemcachedTestSuite) TestSaveThrowError() {
	memcached := memcache.New("127.0.0.1:22222")

	cache := NewMemcached(memcached)

	s.Assert().Error(cache.Save("foo", "bar", 0))
}

func (s *MemcachedTestSuite) TestSave() {
	s.Assert().Nil(s.cache.Save("foo", "bar", 0))
}

func (s *MemcachedTestSuite) TestFetchThrowError() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	memcached := memcache.New("127.0.0.1:22222")
	cache := NewMemcached(memcached)

	result, err := cache.Fetch(key)

	s.Assert().Error(err)
	s.Assert().Empty(result)
}

func (s *MemcachedTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *MemcachedTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().True(s.cache.Contains("foo"))
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *MemcachedTestSuite) TestDeleteThrowError() {
	s.Assert().Error(s.cache.Delete("bar"))
}

func (s *MemcachedTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Delete("foo"))
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *MemcachedTestSuite) TestFlushThrowError() {
	memcached := memcache.New("127.0.0.1:22222")

	cache := NewMemcached(memcached)

	s.Assert().Error(cache.Flush())
}

func (s *MemcachedTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Flush())
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *MemcachedTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	cache := NewMemcached(memcache.New("127.0.0.1:22222"))

	result := cache.FetchMulti([]string{"foo"})

	s.Assert().Len(result, 0)
}

func (s *MemcachedTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.Assert().Len(result, 2)
}

func (s *MemcachedTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.Assert().Len(result, 1)
}

func TestMemcachedRunSuite(t *testing.T) {
	suite.Run(t, new(MemcachedTestSuite))
}
