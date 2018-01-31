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

// NewChain - Create an instance of Chain
func NewChain(drivers ...Cache) *Chain {
    return &Chain{drivers}
}

// Check if cached key exists in one of cache storage
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

// Retrieves the value of the first cache storage found
func (c *Chain) Fetch(key string) (string, error) {

    errs := []string{}

    for _, driver := range c.drivers {
        if value, err := driver.Fetch(key); err == nil {
            return value, nil
        } else {
            errs = append(errs, err.Error())
        }
    }

    return "", fmt.Errorf("Key not found in cache chain. Errors: %s", strings.Join(errs, ","))
}

// Retrieve multiple cached value from keys of the first cache storage found
func (c *Chain) FetchMulti(keys []string) map[string]string {
    result := make(map[string]string)

    for _, key := range keys {
        if value, err := c.Fetch(key); err == nil {
            result[key] = value
        }
    }

    return result
}

// Remove all cached keys of the all cache storages
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
