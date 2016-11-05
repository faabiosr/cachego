package cachego

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
	"net"
	"testing"
	"time"
)

type MongoTestSuite struct {
	suite.Suite

	assert  *assert.Assertions
	cache   Cache
	session *mgo.Session
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

	s.cache = &Mongo{
		session.DB("cache").C("cache"),
	}

	s.assert = assert.New(s.T())
}

func (s *MongoTestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 10))
}

func (s *MongoTestSuite) TestFetchThrowError() {
	result, err := s.cache.Fetch("bar")

	s.assert.Error(err)
	s.assert.Empty(result)
}

func (s *MongoTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.assert.EqualError(err, "Cache expired")
	s.assert.Empty(result)
}

func (s *MongoTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *MongoTestSuite) TestFetchLongCacheDuration() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 10*time.Second)
	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *MongoTestSuite) TestContains() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *MongoTestSuite) TestDeleteThrowError() {
	s.assert.Error(s.cache.Delete("bar"))
}

func (s *MongoTestSuite) TestDelete() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
}

func (s *MongoTestSuite) TestFlush() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *MongoTestSuite) TestFetchMulti() {
	s.cache.Save("foo", "bar", 0)
	s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *MongoTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestMongoRunSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}
