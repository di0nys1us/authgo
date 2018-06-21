package sqlgo

import (
	"database/sql"
	"io"
	"log"
)

// DB TODO
type DB struct {
	*sql.DB
}

// Start TODO
func (db *DB) Start() (*Tx, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	return &Tx{Tx: tx}, nil
}

// NewDB TODO
func NewDB() (*DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres host=artemis sslmode=disable")

	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

// Tx TODO
type Tx struct {
	*sql.Tx
}

// Scanner TODO
type Scanner interface {
	Scan(dest ...interface{}) error
}

// Close TODO
func Close(c io.Closer) {
	err := c.Close()

	if err != nil {
		log.Println(err)
	}
}

// QueryerRow TODO
type QueryerRow interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}
