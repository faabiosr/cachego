package cachego

import "fmt"

type err string

// Error returns the string error value.
func (e err) Error() string {
	return string(e)
}

const (
	// ErrCacheExpired returns an error when the cache key was expired.
	ErrCacheExpired = err("cache expired")

	// ErrFlush returns an error when flush fails.
	ErrFlush = err("unable to flush")

	// ErrSave returns an error when save fails.
	ErrSave = err("unable to save")

	// ErrDelete returns an error when deletion fails.
	ErrDelete = err("unable to delete")

	// ErrDecode returns an errors when decode fails.
	ErrDecode = err("unable to decode")
)

// Wrap returns a new error that adds additional error as a context.
func Wrap(err, additionalErr error) error {
	return fmt.Errorf("%s: %w", additionalErr, err)
}
