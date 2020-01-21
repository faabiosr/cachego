package cachego

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
)

type MongoTestSuite struct {
	suite.Suite

	cache Cache
}

func (s *MongoTestSuite) SetupTest() {
	address := "localhost:27017"

	if _, err := net.Dial("tcp", address); err != nil {
		s.T().Skip()
	}

	session, err := mgo.Dial(address)

	if err != nil {
		s.T().Skip()
	}

	s.cache = NewMongo(session.DB("cache").C("cache"))
}

func (s *MongoTestSuite) TestSave() {
	s.Assert().Nil(s.cache.Save("foo", "bar", 10))
}

func (s *MongoTestSuite) TestFetchThrowError() {
	result, err := s.cache.Fetch("bar")

	s.Assert().Error(err)
	s.Assert().Empty(result)
}

func (s *MongoTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.Assert().EqualError(err, "Cache expired")
	s.Assert().Empty(result)
}

func (s *MongoTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *MongoTestSuite) TestFetchLongCacheDuration() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 10*time.Second)
	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *MongoTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().True(s.cache.Contains("foo"))
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *MongoTestSuite) TestDeleteThrowError() {
	s.Assert().Error(s.cache.Delete("bar"))
}

func (s *MongoTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Delete("foo"))
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *MongoTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Flush())
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *MongoTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.Assert().Len(result, 2)
}

func (s *MongoTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.Assert().Len(result, 1)
}

func TestMongoRunSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}
