package cachego

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type Sqlite3TestSuite struct {
	suite.Suite

	cache Cache
	db    *sql.DB
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

	s.db = db
}

func (s *Sqlite3TestSuite) TearDownTest() {
	os.Remove(dbPath)
}

func (s *Sqlite3TestSuite) TestCreateInstanceThrowAnError() {
	s.db.Close()

	_, err := NewSqlite3(s.db, cacheTable)

	s.Assert().Error(err)
}

func (s *Sqlite3TestSuite) TestSaveThrowAnError() {
	s.db.Close()

	s.Assert().Error(s.cache.Save("foo", "bar", 0))
}

func (s *Sqlite3TestSuite) TestSaveThrowAnErrorWhenDropTable() {
	_, _ = s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.Assert().Error(s.cache.Save("foo", "bar", 0))
}

func (s *Sqlite3TestSuite) TestSave() {
	s.Assert().Nil(s.cache.Save("foo", "bar", 0))
}

func (s *Sqlite3TestSuite) TestFetchThrowAnError() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 1)

	result, err := s.cache.Fetch(key)

	s.Assert().Error(err)
	s.Assert().Empty(result)
}

func (s *Sqlite3TestSuite) TestFetch() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 0)

	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *Sqlite3TestSuite) TestFetchWithLongLifetime() {
	key := "foo"
	value := "bar"

	_ = s.cache.Save(key, value, 10*time.Second)

	result, err := s.cache.Fetch(key)

	s.Assert().Nil(err)
	s.Assert().Equal(value, result)
}

func (s *Sqlite3TestSuite) TestContainsThrowAnError() {
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *Sqlite3TestSuite) TestContains() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().True(s.cache.Contains("foo"))
	s.Assert().False(s.cache.Contains("bar"))
}

func (s *Sqlite3TestSuite) TestDeleteThrowAnError() {
	s.db.Close()

	s.Assert().Error(
		s.cache.Delete("cccc"),
	)
}

func (s *Sqlite3TestSuite) TestDeleteThrowAnErrorWhenDropTable() {
	_, _ = s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.Assert().Error(
		s.cache.Delete("cccc"),
	)
}

func (s *Sqlite3TestSuite) TestDelete() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Delete("foo"))
	s.Assert().False(s.cache.Contains("foo"))
	s.Assert().Nil(s.cache.Delete("foo"))
}

func (s *Sqlite3TestSuite) TestFlushThrowAnError() {
	s.db.Close()

	s.Assert().Error(s.cache.Flush())
}

func (s *Sqlite3TestSuite) TestFlushThrowAnErrorWhenDropTable() {
	_, _ = s.db.Exec(fmt.Sprintf("DROP TABLE %s;", cacheTable))

	s.Assert().Error(s.cache.Flush())
}

func (s *Sqlite3TestSuite) TestFlush() {
	_ = s.cache.Save("foo", "bar", 0)

	s.Assert().Nil(s.cache.Flush())
	s.Assert().False(s.cache.Contains("foo"))
}

func (s *Sqlite3TestSuite) TestFetchMultiReturnNoItemsWhenThrowAnError() {
	s.db.Close()

	result := s.cache.FetchMulti([]string{"foo"})

	s.Assert().Len(result, 0)
}

func (s *Sqlite3TestSuite) TestFetchMulti() {
	_ = s.cache.Save("foo", "bar", 0)
	_ = s.cache.Save("john", "doe", 0)

	result := s.cache.FetchMulti([]string{"foo", "john"})

	s.Assert().Len(result, 2)
}

func (s *Sqlite3TestSuite) TestFetchMultiWhenOnlyOneOfKeysExists() {
	_ = s.cache.Save("foo", "bar", 0)

	result := s.cache.FetchMulti([]string{"foo", "alice"})

	s.Assert().Len(result, 1)
}

func TestSqlite3RunSuite(t *testing.T) {
	suite.Run(t, new(Sqlite3TestSuite))
}
