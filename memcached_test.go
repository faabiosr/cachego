package cachego

import (
	"net"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func TestMemcached(t *testing.T) {
	address := "localhost:11211"
	conn := memcache.New(address)

	if _, err := net.Dial("tcp", address); err != nil {
		t.Skip(err)
	}

	c := NewMemcached(conn)

	if err := c.Save(testKey, testValue, 1*time.Nanosecond); err != nil {
		t.Errorf("save fail: expected nil, got %v", err)
	}

	if _, err := c.Fetch("bar"); err == nil {
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

	c = NewMemcached(
		memcache.New("localhost:22222"),
	)

	if values := c.FetchMulti([]string{testKey, "bar"}); len(values) != 0 {
		t.Errorf("fetch multi failed: expected %d, got %d", 0, len(values))
	}
}
