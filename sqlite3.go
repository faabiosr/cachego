package cachego

import (
	"database/sql"
	"fmt"
	errors "github.com/pkg/errors"
	"time"
)

type (
	// Sqlite3 store for caching data
	Sqlite3 struct {
		db    *sql.DB
		table string
	}
)

// NewSqlite3 - Create an instance of Sqlite3
func NewSqlite3(db *sql.DB, table string) (*Sqlite3, error) {
	if err := createTable(db, table); err != nil {
		return nil, errors.Wrap(err, "Unable to create database table")
	}

	return &Sqlite3{db, table}, nil
}

func createTable(db *sql.DB, table string) error {
	stmt := `CREATE TABLE IF NOT EXISTS %s (
        key text PRIMARY KEY,
        value text NOT NULL,
        lifetime integer NOT NULL
    );`

	_, err := db.Exec(fmt.Sprintf(stmt, table))

	return err
}

// Check if cached key exists in SQL storage
func (s *Sqlite3) Contains(key string) bool {
	if _, err := s.Fetch(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from Sqlite3 storage
func (s *Sqlite3) Delete(key string) error {
	tx, err := s.db.Begin()

	if err != nil {
		return errors.Wrap(err, "Unable to delete")
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("DELETE FROM %s WHERE key = ?", s.table),
	)

	if err != nil {
		return errors.Wrap(err, "Unable to delete")
	}

	defer stmt.Close()

	_, err = stmt.Exec(key)

	if err != nil {
		return errors.Wrap(err, "Unable to delete")
	}

	tx.Commit()

	return nil
}

// Retrieve the cached value from key of the Sqlite3 storage
func (s *Sqlite3) Fetch(key string) (string, error) {
	stmt, err := s.db.Prepare(
		fmt.Sprintf("SELECT value, lifetime FROM %s WHERE key = ?", s.table),
	)

	if err != nil {
		return "", errors.Wrap(err, "Unable to retrieve the value")
	}

	defer stmt.Close()

	var value string
	var lifetime int64

	err = stmt.QueryRow(key).Scan(&value, &lifetime)

	if err != nil {
		return "", errors.Wrap(err, "Unable to retrieve the value")
	}

	if lifetime == 0 {
		return value, nil
	}

	if lifetime <= time.Now().Unix() {
		s.Delete(key)

		return "", errors.New("Cache expired")
	}

	return value, nil
}

// Retrieve multiple cached value from keys of the Sqlite3 storage
func (s *Sqlite3) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := s.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Remove all cached keys in Sqlite3 storage
func (s *Sqlite3) Flush() error {
	tx, err := s.db.Begin()

	if err != nil {
		return errors.Wrap(err, "Unable to flush")
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("DELETE FROM %s", s.table),
	)

	if err != nil {
		return errors.Wrap(err, "Unable to flush")
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		return errors.Wrap(err, "Unable to flush")
	}

	tx.Commit()

	return nil
}

// Save a value in Sqlite3 storage by key
func (s *Sqlite3) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	tx, err := s.db.Begin()

	if err != nil {
		return errors.Wrap(err, "Unable to save")
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT OR REPLACE INTO %s (key, value, lifetime) VALUES (?, ?, ?)", s.table),
	)

	if err != nil {
		return errors.Wrap(err, "Unable to save")
	}

	defer stmt.Close()

	_, err = stmt.Exec(key, value, duration)

	if err != nil {
		return errors.Wrap(err, "Unable to save")
	}

	tx.Commit()

	return nil
}
