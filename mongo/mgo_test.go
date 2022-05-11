package mongo

import (
	"net"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	testKeyMgo   = "foo"
	testValueMgo = "bar"
)

func TestMgo(t *testing.T) {
	address := "localhost:27017"

	if _, err := net.Dial("tcp", address); err != nil {
		t.Skip(err)
	}

	session, err := mgo.Dial(address)
	if err != nil {
		t.Skip(err)
	}

	c := New(session.DB("cache").C("cache"))

	if err := c.Save(testKeyMgo, testValueMgo, 1*time.Nanosecond); err != nil {
		t.Errorf("save fail: expected nil, got %v", err)
	}

	if _, err := c.Fetch(testKeyMgo); err == nil {
		t.Errorf("fetch fail: expected an error, got %v", err)
	}

	_ = c.Save(testKeyMgo, testValueMgo, 10*time.Second)

	if res, _ := c.Fetch(testKeyMgo); res != testValueMgo {
		t.Errorf("fetch fail, wrong value: expected %s, got %s", testValueMgo, res)
	}

	_ = c.Save(testKeyMgo, testValueMgo, 0)

	if !c.Contains(testKeyMgo) {
		t.Errorf("contains failed: the key %s should be exist", testKeyMgo)
	}

	_ = c.Save("bar", testValueMgo, 0)

	if values := c.FetchMulti([]string{testKeyMgo, "bar"}); len(values) == 0 {
		t.Errorf("fetch multi failed: expected %d, got %d", 2, len(values))
	}

	if err := c.Flush(); err != nil {
		t.Errorf("flush failed: expected nil, got %v", err)
	}

	if c.Contains(testKeyMgo) {
		t.Errorf("contains failed: the key %s should not be exist", testKeyMgo)
	}
}
