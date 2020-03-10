package cachego

import (
	"database/sql"
	"fmt"
	"time"
)

type (
	sqlite3 struct {
		db    *sql.DB
		table string
	}
)

// NewSqlite3 creates an instance of Sqlite3 cache driver
func NewSqlite3(db *sql.DB, table string) (Cache, error) {
	if err := createTable(db, table); err != nil {
		return nil, err
	}

	return &sqlite3{db, table}, nil
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

// Contains checks if cached key exists in Sqlite3 storage
func (s *sqlite3) Contains(key string) bool {
	if _, err := s.Fetch(key); err != nil {
		return false
	}

	return true
}

// Delete the cached key from Sqlite3 storage
func (s *sqlite3) Delete(key string) error {
	tx, err := s.db.Begin()

	if err != nil {
		return Wrap(ErrDelete, err)
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("DELETE FROM %s WHERE key = ?", s.table),
	)

	if err != nil {
		return Wrap(ErrDelete, err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(key)

	if err != nil {
		return Wrap(ErrDelete, err)
	}

	_ = tx.Commit()

	return nil
}

// Fetch retrieves the cached value from key of the Sqlite3 storage
func (s *sqlite3) Fetch(key string) (string, error) {
	stmt, err := s.db.Prepare(
		fmt.Sprintf("SELECT value, lifetime FROM %s WHERE key = ?", s.table),
	)

	if err != nil {
		return "", err
	}

	defer func() {
		_ = stmt.Close()
	}()

	var value string
	var lifetime int64

	err = stmt.QueryRow(key).Scan(&value, &lifetime)

	if err != nil {
		return "", err
	}

	if lifetime == 0 {
		return value, nil
	}

	if lifetime <= time.Now().Unix() {
		_ = s.Delete(key)

		return "", ErrCacheExpired
	}

	return value, nil
}

// FetchMulti retrieves multiple cached value from keys of the Sqlite3 storage
func (s *sqlite3) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := s.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

// Flush removes all cached keys of the Sqlite3 storage
func (s *sqlite3) Flush() error {
	tx, err := s.db.Begin()

	if err != nil {
		return Wrap(ErrFlush, err)
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("DELETE FROM %s", s.table),
	)

	if err != nil {
		return Wrap(ErrFlush, err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec()

	if err != nil {
		return Wrap(ErrFlush, err)
	}

	_ = tx.Commit()

	return nil
}

// Save a value in Sqlite3 storage by key
func (s *sqlite3) Save(key string, value string, lifeTime time.Duration) error {
	duration := int64(0)

	if lifeTime > 0 {
		duration = time.Now().Unix() + int64(lifeTime.Seconds())
	}

	tx, err := s.db.Begin()

	if err != nil {
		return Wrap(ErrSave, err)
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT OR REPLACE INTO %s (key, value, lifetime) VALUES (?, ?, ?)", s.table),
	)

	if err != nil {
		return Wrap(ErrSave, err)
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(key, value, duration)

	if err != nil {
		return Wrap(ErrSave, err)
	}

	_ = tx.Commit()

	return nil
}
