package cachego

import (
	"fmt"
	"os"
	"testing"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/stretchr/testify/suite"
)

type BoltTestSuite struct {
	suite.Suite

	cache     Cache
	db        *bolt.DB
	directory string
}

func (s *BoltTestSuite) SetupTest() {
	s.directory = "./cache-dir/"

	_ = os.Mkdir(s.directory, 0777)

	db, err := bolt.Open(s.directory+"cachego.db", 0600, nil)

	if err != nil {
		s.T().Skip()
	}

	s.db = db
	s.cache = NewBolt(s.db)
}

func (s *BoltTestSuite) TearDownTest() {
	s.db.Close()
}

func (s *BoltTestSuite) TestSave() {
	s.Assert().Nil(s.cache.Save("foo", "bar", 0))
}

func (s *BoltTestSuite) TestSaveThrowError() {
	s.db.Close()

	opts := &bolt.Options{ReadOnly: true}
	db, err := bolt.Open(s.directory+"cachego.db", 0666, opts)

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	cache := NewBolt(db)
	err = cache.Save("foo", "bar", 0)

	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), ErrBoltSave)
}

func (s *BoltTestSuite) TestFetchThrowErrorWhenBucketNotFound() {
	s.cache.Flush()

	result, err := s.cache.Fetch("foo")

	s.Assert().Empty(result)
	s.Assert().EqualError(err, ErrBoltBucketNotFound.Error())
}

func (s *BoltTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.Assert().Empty(result)
	s.Assert().EqualError(err, ErrBoltCacheExpired.Error())
}

func (s *BoltTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)
	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *BoltTestSuite) TestFetchLongCacheDuration() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 10*time.Second)
	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *BoltTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().True(s.cache.Contains("foo"))
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *BoltTestSuite) TestDeleteThrowErrorWhenBucketNotFound() {
	s.cache.Flush()

	err := s.cache.Delete("foo")

	s.Assert().EqualError(err, ErrBoltBucketNotFound.Error())
}

func (s *BoltTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Delete("foo"))
	s.Assert().False(s.cache.Contains("foo"))
	s.Assert().Nil(s.cache.Delete("bar"))
}

func (s *BoltTestSuite) TestFlushThrowErrorWhenBucketNotFound() {
	err := s.cache.Flush()

	s.Assert().Error(err)
	s.Assert().Contains(err.Error(), ErrBoltFlush)
}

func (s *BoltTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Flush())
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *BoltTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	s.cache.Flush()
	result := s.cache.FetchMulti([]string{"foo"})

	s.Assert().Len(result, 0)
}

func (s *BoltTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.Assert().Len(result, 2)
}

func (s *BoltTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.Assert().Len(result, 1)
}

func TestBoltRunSuite(t *testing.T) {
	suite.Run(t, new(BoltTestSuite))
}
