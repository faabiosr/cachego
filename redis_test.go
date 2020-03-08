package cachego

import (
	"net"
	"testing"

	"github.com/stretchr/testify/suite"
	rd "gopkg.in/redis.v4"
)

type RedisTestSuite struct {
	suite.Suite

	cache Cache
}

func (s *RedisTestSuite) SetupTest() {
	conn := rd.NewClient(&rd.Options{
		Addr: ":6379",
	})

	if _, err := net.Dial("tcp", "localhost:6379"); err != nil {
		s.T().Skip()
	}

	s.cache = NewRedis(conn)
}

func (s *RedisTestSuite) TestSaveThrowError() {
	conn := rd.NewClient(&rd.Options{
		Addr: ":6380",
	})

	cache := NewRedis(conn)

	s.Assert().Error(cache.Save("foo", "bar", 0))
}

func (s *RedisTestSuite) TestSave() {
	s.Assert().Nil(s.cache.Save("foo", "bar", 0))
}

func (s *RedisTestSuite) TestFetchThrowError() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	conn := rd.NewClient(&rd.Options{
		Addr: ":6380",
	})

	cache := NewRedis(conn)

	result, err := cache.Fetch(key)

	s.Assert().Error(err)
	s.Assert().Empty(result)
}

func (s *RedisTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *RedisTestSuite) TestContainsThrowError() {
	conn := rd.NewClient(&rd.Options{
		Addr: ":6380",
	})

	cache := NewRedis(conn)

	s.Assert().False(cache.Contains("bar"))
}

func (s *RedisTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().True(s.cache.Contains("foo"))
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *RedisTestSuite) TestDeleteThrowError() {
	conn := rd.NewClient(&rd.Options{
		Addr: ":6380",
	})

	cache := NewRedis(conn)
	s.Assert().Error(cache.Delete("bar"))
}

func (s *RedisTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Delete("foo"))
	s.Assert().False(s.cache.Contains("foo"))
	s.Assert().Nil(s.cache.Delete("foo"))
}

func (s *RedisTestSuite) TestFlushThrowError() {
	conn := rd.NewClient(&rd.Options{
		Addr: ":6380",
	})

	cache := NewRedis(conn)

	s.Assert().Error(cache.Flush())
}

func (s *RedisTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Flush())
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *RedisTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	cache := NewRedis(rd.NewClient(&rd.Options{
		Addr: ":6380",
	}))

	result := cache.FetchMulti([]string{"foo"})

	s.Assert().Len(result, 0)
}

func (s *RedisTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.Assert().Len(result, 2)
}

func (s *RedisTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.Assert().Len(result, 1)
}

func TestRedisRunSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
