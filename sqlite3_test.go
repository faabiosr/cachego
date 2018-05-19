package cachego

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type Sqlite3TestSuite struct {
	suite.Suite

	assert *assert.Assertions
	cache  Cache
	db     *sql.DB
}

var (
	cacheTable = "cache"
	dbPath     = "./cache.db"
)

func (s *Sqlite3TestSuite) SetupTest() {

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		s.T().Skip()
	}

	s.cache, err = NewSqlite3(db, cacheTable)

	if err != nil {
		s.T().Skip()
	}

	s.assert = assert.New(s.T())
	s.db = db
}

func (s *Sqlite3TestSuite) TearDownTest() {
	os.Remove(dbPath)
}

func (s *Sqlite3TestSuite) TestCreateInstanceThrowAnError() {
	s.db.Close()

	_, err := NewSqlite3(s.db, cacheTable)

	s.assert.Error(err)
}

func (s *Sqlite3TestSuite) TestSaveThrowAnError() {
	s.db.Close()

	s.assert.Error(s.cache.Save("foo", "bar", 0))
}

func (s *Sqlite3TestSuite) TestSaveThrowAnErrorWhenDropTable() {
	s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.assert.Error(s.cache.Save("foo", "bar", 0))
}

func (s *Sqlite3TestSuite) TestSave() {
	s.assert.Nil(s.cache.Save("foo", "bar", 0))
}

func (s *Sqlite3TestSuite) TestFetchThrowAnError() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 1)

	result, err := s.cache.Fetch(key)

	s.assert.Error(err)
	s.assert.Empty(result)
}

func (s *Sqlite3TestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *Sqlite3TestSuite) TestFetchWithLongLifetime() {
	key := "foo"
	value := "bar"

	s.cache.Save(key, value, 10*time.Second)

	result, err := s.cache.Fetch(key)

	s.assert.Nil(err)
	s.assert.Equal(value, result)
}

func (s *Sqlite3TestSuite) TestContainsThrowAnError() {
	s.assert.False(s.cache.Contains("bar"))
}

func (s *Sqlite3TestSuite) TestContains() {
	s.cache.Save("foo", "bar", 0)

	s.assert.True(s.cache.Contains("foo"))
	s.assert.False(s.cache.Contains("bar"))
}

func (s *Sqlite3TestSuite) TestDeleteThrowAnError() {
	s.db.Close()

	s.assert.Error(
		s.cache.Delete("cccc"),
	)
}

func (s *Sqlite3TestSuite) TestDeleteThrowAnErrorWhenDropTable() {
	s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.assert.Error(
		s.cache.Delete("cccc"),
	)
}

func (s *Sqlite3TestSuite) TestDelete() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Delete("foo"))
	s.assert.False(s.cache.Contains("foo"))
	s.assert.Nil(s.cache.Delete("foo"))
}

func (s *Sqlite3TestSuite) TestFlushThrowAnError() {
	s.db.Close()

	s.assert.Error(s.cache.Flush())
}

func (s *Sqlite3TestSuite) TestFlushThrowAnErrorWhenDropTable() {
	s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.assert.Error(s.cache.Flush())
}

func (s *Sqlite3TestSuite) TestFlush() {
	s.cache.Save("foo", "bar", 0)

	s.assert.Nil(s.cache.Flush())
	s.assert.False(s.cache.Contains("foo"))
}

func (s *Sqlite3TestSuite) TestFetchMultiReturnNoItemsWhenThrowAnError() {
	s.db.Close()

	result := s.cache.FetchMulti([]string{"foo"})

	s.assert.Len(result, 0)
}

func (s *Sqlite3TestSuite) TestFetchMulti() {
	s.cache.Save("foo", "bar", 0)
	s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.assert.Len(result, 2)
}

func (s *Sqlite3TestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.assert.Len(result, 1)
}

func TestSqlite3RunSuite(t *testing.T) {
	suite.Run(t, new(Sqlite3TestSuite))
}
