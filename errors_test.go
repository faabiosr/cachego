package cachego

import (
	"errors"
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

func TestWrap(t *testing.T) {
	additionalErr := errors.New("failed")
	err := Wrap(ErrSave, additionalErr)

	if !errors.Is(err, ErrSave) {
		t.Errorf("wrap failed: expected true")
	}
}
