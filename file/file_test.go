package file

import (
	"os"
	"testing"
	"time"
)

const (
	testKey   = "foo"
	testValue = "bar"
)

func TestFile(t *testing.T) {
	dir, err := os.MkdirTemp("", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	c := New(dir)

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

	c = New("./test/")

	if err := c.Save(testKey, testValue, 0); err == nil {
		t.Errorf("save failed: expected an error, got %v", err)
	}

	if _, err := c.Fetch(testKey); err == nil {
		t.Errorf("fetch failed: expected and error, got %v", err)
	}

	if err := c.Flush(); err == nil {
		t.Errorf("flush failed: expected an error, got %v", err)
	}
}
