package sqlite3

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	testKey    = "foo"
	testValue  = "bar"
	testTable  = "cache"
	testDBPath = "/cache.db"
)

func TestSqlite3(t *testing.T) {
	dir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite3", dir+testDBPath)
	if err != nil {
		t.Skip(err)
	}

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	c, err := New(db, testTable)
	if err != nil {
		t.Skip(err)
	}

	if err := c.Save(testKey, testValue, 1*time.Nanosecond); err != nil {
		t.Errorf("save fail: expected nil, got %v", err)
	}

	if _, err := c.Fetch(testKey); err == nil {
		t.Errorf("fetch fail: expected an error, got %v", err)
	}

	_ = c.Save(testKey, testValue, 10*time.Second)

	if res, _ := c.Fetch(testKey); res != testValue {
		t.Errorf("fetch fail, wrong value: expected %s, got %s", testValue, res)
	}

	_ = c.Save(testKey, testValue, 0)

	if !c.Contains(testKey) {
		t.Errorf("contains failed: the key %s should be exist", testKey)
	}

	_ = c.Save("bar", testValue, 0)

	if values := c.FetchMulti([]string{testKey, "bar"}); len(values) == 0 {
		t.Errorf("fetch multi failed: expected %d, got %d", 2, len(values))
	}

	if err := c.Flush(); err != nil {
		t.Errorf("flush failed: expected nil, got %v", err)
	}

	if c.Contains(testKey) {
		t.Errorf("contains failed: the key %s should not be exist", testKey)
	}
}

func TestSqlite3Fail(t *testing.T) {
	dir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	db, _ := sql.Open("sqlite3", dir+testDBPath)
	_ = db.Close()

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	if _, err := New(db, testTable); err == nil {
		t.Errorf("constructor failed: expected an error, got %v", err)
	}

	db, _ = sql.Open("sqlite3", testDBPath)
	c, _ := New(db, testTable)
	_ = db.Close()

	if err := c.Save(testKey, testValue, 0); err == nil {
		t.Errorf("save failed: expected an error, got %v", err)
	}

	if err := c.Delete(testKey); err == nil {
		t.Errorf("delete failed: expected an error, got %v", err)
	}

	if err := c.Flush(); err == nil {
		t.Errorf("flush failed: expected an error, got %v", err)
	}

	db, _ = sql.Open("sqlite3", testDBPath)
	c, _ = New(db, testTable)

	_, _ = db.Exec(fmt.Sprintf("DROP TABLE %s;", testTable))

	if err := c.Save(testKey, testValue, 0); err == nil {
		t.Errorf("save failed: expected an error, got %v", err)
	}

	if err := c.Delete(testKey); err == nil {
		t.Errorf("delete failed: expected an error, got %v", err)
	}

	if err := c.Flush(); err == nil {
		t.Errorf("flush failed: expected an error, got %v", err)
	}
}
