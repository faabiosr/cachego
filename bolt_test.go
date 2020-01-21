package cachego

import (
	"fmt"
	"os"
	"testing"
	"time"

	bolt "github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BoltTestSuite struct {
	suite.Suite

	assert    *assert.Assertions
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
	s.assert = assert.New(s.T())
}

func (s *BoltTestSuite) TearDownTest() {
	s.db.Close()
}

func (s *BoltTestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 0))
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

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrBoltSave)
}

func (s *BoltTestSuite) TestFetchThrowErrorWhenBucketNotFound() {
	s.cache.Flush()

	result, err := s.cache.Fetch("foo")

	s.assert.Empty(result)
	s.assert.EqualError(err, ErrBoltBucketNotFound.Error())
}

func (s *BoltTestSuite) TestFetchThrowErrorWhenExpired() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1*time.Second)

	time.Sleep(1 * time.Second)

	result, err := s.cache.Fetch(key)

	s.assert.Empty(result)
	s.assert.EqualError(err, ErrBoltCacheExpired.Error())
}

func (s *BoltTestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)
	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *BoltTestSuite) TestFetchLongCacheDuration() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 10*time.Second)
	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *BoltTestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *BoltTestSuite) TestDeleteThrowErrorWhenBucketNotFound() {
	s.cache.Flush()

	err := s.cache.Delete("foo")

	s.assert.EqualError(err, ErrBoltBucketNotFound.Error())
}

func (s *BoltTestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
	s.assert.Nil(s.cache.Delete("bar"))
}

func (s *BoltTestSuite) TestFlushThrowErrorWhenBucketNotFound() {
	err := s.cache.Flush()

	s.assert.Error(err)
	s.assert.Contains(err.Error(), ErrBoltFlush)
}

func (s *BoltTestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *BoltTestSuite) TestFetchMultiReturnNoItemsWhenThrowError() {
	s.cache.Flush()
	result := s.cache.FetchMulti([]string{"foo"})

	s.assert.Len(result, 0)
}

func (s *BoltTestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *BoltTestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestBoltRunSuite(t *testing.T) {
	suite.Run(t, new(BoltTestSuite))
}
