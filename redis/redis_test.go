package redis

import (
	"net"
	"testing"
	"time"

	rd "github.com/go-redis/redis/v8"
)

const (
	testKey   = "foo"
	testValue = "bar"
)

func TestRedis(t *testing.T) {
	conn := rd.NewClient(&rd.Options{
		Addr: ":6379",
	})

	if _, err := net.Dial("tcp", "localhost:6379"); err != nil {
		t.Skip(err)
	}

	c := New(conn)

	if err := c.Save(testKey, testValue, 10*time.Second); err != nil {
		t.Errorf("save fail: expected nil, got %v", err)
	}

	if res, _ := c.Fetch(testKey); res != testValue {
		t.Errorf("fetch fail, wrong value: expected %s, got %s", testValue, res)
	}

	if _, err := c.Fetch("bar"); err == nil {
		t.Errorf("fetch fail: expected an error, got %v", err)
	}

	if !c.Contains(testKey) {
		t.Errorf("contains failed: the key %s should be exist", testKey)
	}

	_ = c.Save("bar", testValue, 0)

	if values := c.FetchMulti([]string{testKey, "bar"}); len(values) == 0 {
		t.Errorf("fetch multi failed: expected %d, got %d", 2, len(values))
	}

	if err := c.Delete(testKey); err != nil {
		t.Errorf("delete failed: expected nil, got %v", err)
	}

	if err := c.Flush(); err != nil {
		t.Errorf("flush failed: expected nil, got %v", err)
	}

	if c.Contains(testKey) {
		t.Errorf("contains failed: the key %s should not be exist", testKey)
	}

	conn = rd.NewClient(&rd.Options{Addr: ":6380"})

	c = New(conn)

	if c.Contains(testKey) {
		t.Errorf("contains failed: the key %s should not be exist", testKey)
	}

	if values := c.FetchMulti([]string{testKey, "bar"}); len(values) != 0 {
		t.Errorf("fetch multi failed: expected %d, got %d", 0, len(values))
	}
}
