package cachego

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	expect := "failed"
	er := err(expect)

	if r := fmt.Sprint(er); r != expect {
		t.Errorf("invalid string: expect %s, got %s", expect, r)
	}
}
