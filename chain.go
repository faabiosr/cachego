package cachego

import (
	"fmt"
	"strings"
	"time"
)

type (
	// Chain storage for dealing with multiple cache storage in the same time
	Chain struct {
		drivers []Cache
	}
)

// NewChain creates an instance of Chain cache driver
func NewChain(drivers ...Cache) Cache {
	return &Chain{drivers}
}

// Contains checks if the cached key exists in one of the cache storages
func (c *Chain) Contains(key string) bool {
	for _, driver := range c.drivers {
		if driver.Contains(key) {
			return true
		}
	}

	return false
}

// Delete the cached key in all cache storages
func (c *Chain) Delete(key string) error {
	for _, driver := range c.drivers {
		if err := driver.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

// Fetch retrieves the value of one of the registred cache storages
func (c *Chain) Fetch(key string) (string, error) {

	errs := []string{}

	for _, driver := range c.drivers {
		value, err := driver.Fetch(key)

		if err == nil {
			return value, nil
		}

		errs = append(errs, err.Error())
	}

	return "", fmt.Errorf("Key not found in cache chain. Errors: %s", strings.Join(errs, ","))
}

// FetchMulti retrieves multiple cached values from one of the registred cache storages
func (c *Chain) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := c.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the registered cache storages
func (c *Chain) Flush() error {
	for _, driver := range c.drivers {
		if err := driver.Flush(); err != nil {
			return err
		}
	}

	return nil
}

// Save a value in all cache storages by key
func (c *Chain) Save(key string, value string, lifeTime time.Duration) error {

	for _, driver := range c.drivers {
		if err := driver.Save(key, value, lifeTime); err != nil {
			return err
		}
	}

	return nil
}
