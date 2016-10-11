package cachego

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/redis.v4"
	"net"
	"testing"
)

type RedisTestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
}

func (s *RedisTestSuite) SetupTest() {
	conn := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	if _, err := net.Dial("tcp", "localhost:6379"); err != nil {
		s.T().Skip()
	}

	s.cache = &Redis{conn}
	s.assert = assert.New(s.T())
}

func (s *RedisTestSuite) TestSaveReturnFalseWhenThrowError() {
	redis := redis.NewClient(&redis.Options{
		Addr: ":6380",
	})

	cache := &Redis{redis}

	s.assert.False(cache.Save("foo", "bar", 0))
}

func (s *RedisTestSuite) TestSave() {
	s.assert.True(s.cache.Save("foo", "bar", 0))
}

func (s *RedisTestSuite) TestFetchReturnFalseWhenThrowError() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	redis := redis.NewClient(&redis.Options{
		Addr: ":6380",
	})
	cache := &Redis{redis}

	result, status := cache.Fetch(key)

	s.assert.False(status)
	s.assert.Empty(result)
}

func (s *RedisTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	result, status := s.cache.Fetch(key)

	s.assert.True(status)
	s.assert.Equal(value, result)
}

func (s *RedisTestSuite) TestContainsReturnFalseWhenThrowError() {
	redis := redis.NewClient(&redis.Options{
		Addr: ":6380",
	})

	cache := &Redis{redis}

	s.assert.False(cache.Contains("bar"))
}

func (s *RedisTestSuite) TestContains() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *RedisTestSuite) TestDeleteReturnFalseWhenThrowError() {
	redis := redis.NewClient(&redis.Options{
		Addr: ":6380",
	})

	cache := &Redis{redis}
	s.assert.False(cache.Delete("bar"))
}

func (s *RedisTestSuite) TestDelete() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
	s.assert.False(s.cache.Delete("foo"))
}

func (s *RedisTestSuite) TestFlushReturnFalseWhenThrowError() {
	redis := redis.NewClient(&redis.Options{
		Addr: ":6380",
	})

	cache := &Redis{redis}

	s.assert.False(cache.Flush())
}

func (s *RedisTestSuite) TestFlush() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *RedisTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	cache := &Redis{
		redis.NewClient(&redis.Options{
			Addr: ":6380",
		}),
	}

	result := cache.FetchMulti([]string{"foo"})

	s.assert.Len(result, 0)
}

func (s *RedisTestSuite) TestFetchMulti() {
	s.cache.Save("foo", "bar", 0)
	s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *RedisTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestRedisRunSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
