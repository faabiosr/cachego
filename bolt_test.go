package cachego

import (
	"fmt"
	"os"
	"testing"
	"time"

	bt "github.com/coreos/bbolt"
)

func TestBolt(t *testing.T) {
	dir := "./cache-dir"
	_ = os.Mkdir(dir, 0777)

	db, err := bt.Open(fmt.Sprintf("%s/cachego.db", dir), 0600, nil)

	if err != nil {
		t.Skip(err)
	}

	defer func() {
		_ = db.Close()
	}()

	c := NewBolt(db)

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

	if err := c.Flush(); err == nil {
		t.Errorf("flush failed: expected error, got %v", err)
	}

	if err := c.Delete(testKey); err == nil {
		t.Errorf("delete failed: expected error, got %v", err)
	}

	if c.Contains(testKey) {
		t.Errorf("contains failed: the key %s should not be exist", testKey)
	}
}

func TestBoltSaveWithReadOnlyDB(t *testing.T) {
	dir := "./cache-dir"
	_ = os.Mkdir(dir, 0777)

	db, err := bt.Open(fmt.Sprintf("%s/cachego.db", dir), 0666, &bt.Options{ReadOnly: true})

	if err != nil {
		t.Skip(err)
	}

	defer func() {
		_ = db.Close()
	}()

	c := NewBolt(db)

	if err := c.Save(testKey, testValue, 0); err == nil {
		t.Errorf("save failed: expected error, got %v", err)

	}
}
