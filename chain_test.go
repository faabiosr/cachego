package cachego

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ChainTestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
}

func (s *ChainTestSuite) SetupTest() {
	s.cache = NewChain(NewMap())

	s.assert = assert.New(s.T())
}

func (s *ChainTestSuite) TestSaveThrowErrorWhenOneOfDriverFail() {
	cache := NewChain(
		NewMap(),
		NewMemcached(memcache.New("127.0.0.1:22222")),
	)

	s.assert.Error(cache.Save("foo", "bar", 0))
}

func (s *ChainTestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 0))
}

func (s *ChainTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.assert.Regexp("^Key not found in cache chain", err)
	s.assert.Empty(result)
}

func (s *ChainTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *ChainTestSuite) TestContains() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *ChainTestSuite) TestDeleteThrowErrorWhenOneOfDriverFail() {
	cache := NewChain(
		NewMap(),
		NewMemcached(memcache.New("127.0.0.1:22222")),
	)

	s.assert.Error(cache.Delete("foo"))
}

func (s *ChainTestSuite) TestDelete() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
}

func (s *ChainTestSuite) TestFlushThrowErrorWhenOneOfDriverFail() {
	cache := NewChain(
		NewMap(),
		NewMemcached(memcache.New("127.0.0.1:22222")),
	)

	s.assert.Error(cache.Flush())
}

func (s *ChainTestSuite) TestFlush() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *ChainTestSuite) TestFetchMulti() {
	s.cache.Save("foo", "bar", 0)
	s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *ChainTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestChainRunSuite(t *testing.T) {
	suite.Run(t, new(ChainTestSuite))
}
