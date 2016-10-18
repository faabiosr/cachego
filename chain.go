package cachego

import (
	"fmt"
	"strings"
	"time"
)

type Chain struct {
	drivers []Cache
}

func (c *Chain) Contains(key string) bool {
	for _, driver := range c.drivers {
		if driver.Contains(key) {
			return true
		}
	}

	return false
}

func (c *Chain) Delete(key string) error {
	for _, driver := range c.drivers {
		if err := driver.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

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

func (c *Chain) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := c.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

func (c *Chain) Flush() error {
	for _, driver := range c.drivers {
		if err := driver.Flush(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Chain) Save(key string, value string, lifeTime time.Duration) error {

	for _, driver := range c.drivers {
		if err := driver.Save(key, value, lifeTime); err != nil {
			return err
		}
	}

	return nil
}
