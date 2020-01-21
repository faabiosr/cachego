package cachego

import (
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/stretchr/testify/suite"
)

type ChainTestSuite struct {
	suite.Suite

	cache Cache
}

func (s *ChainTestSuite) SetupTest() {
	s.cache = NewChain(NewMap())
}

func (s *ChainTestSuite) TestSaveThrowErrorWhenOneOfDriverFail() {
	cache := NewChain(
		NewMap(),
		NewMemcached(memcache.New("127.0.0.1:22222")),
	)

	s.Assert().Error(cache.Save("foo", "bar", 0))
}

func (s *ChainTestSuite) TestSave() {
	s.Assert().Nil(s.cache.Save("foo", "bar", 0))
}

func (s *ChainTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.Assert().Regexp("^Key not found in cache chain", err)
	s.Assert().Empty(result)
}

func (s *ChainTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *ChainTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().True(s.cache.Contains("foo"))
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *ChainTestSuite) TestDeleteThrowErrorWhenOneOfDriverFail() {
	cache := NewChain(
		NewMap(),
		NewMemcached(memcache.New("127.0.0.1:22222")),
	)

	s.Assert().Error(cache.Delete("foo"))
}

func (s *ChainTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Delete("foo"))
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *ChainTestSuite) TestFlushThrowErrorWhenOneOfDriverFail() {
	cache := NewChain(
		NewMap(),
		NewMemcached(memcache.New("127.0.0.1:22222")),
	)

	s.Assert().Error(cache.Flush())
}

func (s *ChainTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Flush())
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *ChainTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.Assert().Len(result, 2)
}

func (s *ChainTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.Assert().Len(result, 1)
}

func TestChainRunSuite(t *testing.T) {
	suite.Run(t, new(ChainTestSuite))
}
