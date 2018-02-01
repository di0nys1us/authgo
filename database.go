package authgo

import (
	"github.com/jmoiron/sqlx"
)

type db struct {
	*sqlx.DB
}

func (db *db) begin() (*tx, error) {
	wrapped, err := db.DB.Beginx()

	if err != nil {
		return nil, err
	}

	return &tx{wrapped}, nil
}

func newDB() (*db, error) {
	wrapped, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		return nil, err
	}

	return &db{wrapped}, nil
}

type tx struct {
	*sqlx.Tx
}

type saver interface {
	save(tx *tx) error
}

type updater interface {
	update(tx *tx) error
}

type deleter interface {
	delete(tx *tx) error
}
