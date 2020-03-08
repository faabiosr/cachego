package cachego

import (
	"time"
)

type (
	chain struct {
		drivers []Cache
	}
)

// NewChain creates an instance of Chain cache driver
func NewChain(drivers ...Cache) Cache {
	return &chain{drivers}
}

// Contains checks if the cached key exists in one of the cache storages
func (c *chain) Contains(key string) bool {
	for _, driver := range c.drivers {
		if driver.Contains(key) {
			return true
		}
	}

	return false
}

// Delete the cached key in all cache storages
func (c *chain) Delete(key string) error {
	for _, driver := range c.drivers {
		if err := driver.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

// Fetch retrieves the value of one of the registred cache storages
func (c *chain) Fetch(key string) (string, error) {
	for _, driver := range c.drivers {
		value, err := driver.Fetch(key)

		if err == nil {
			return value, nil
		}
	}

	return "", err("key not found in cache chain")
}

// FetchMulti retrieves multiple cached values from one of the registred cache storages
func (c *chain) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := c.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the registered cache storages
func (c *chain) Flush() error {
	for _, driver := range c.drivers {
		if err := driver.Flush(); err != nil {
			return err
		}
	}

	return nil
}

// Save a value in all cache storages by key
func (c *chain) Save(key string, value string, lifeTime time.Duration) error {
	for _, driver := range c.drivers {
		if err := driver.Save(key, value, lifeTime); err != nil {
			return err
		}
	}

	return nil
}
