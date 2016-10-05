package cachego

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/redis.v4"
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

	s.cache = &Redis{conn}
	s.assert = assert.New(s.T())
}

func (s *RedisTestSuite) TestSaveReturnFalseWhenThrowError() {
	redis := redis.NewClient(&redis.Options{
		Addr: ":6380",
	})

	cache := &Redis{redis}

	s.assert.False(cache.Save("foo", "bar"))
}

func (s *RedisTestSuite) TestSave() {
	s.assert.True(s.cache.Save("foo", "bar"))
}

func (s *RedisTestSuite) TestFetchReturnFalseWhenThrowError() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value)

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

	s.cache.Save(key, value)

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
	s.cache.Save("foo", "bar")

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
	s.cache.Save("foo", "bar")

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
	s.cache.Save("foo", "bar")

	s.assert.True(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func TestRedisRunSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
