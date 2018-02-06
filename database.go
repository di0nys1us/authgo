package authgo

import (
	"log"

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

func (db *db) do(fn func(tx *tx)) error {
	txx, err := db.begin()

	if err != nil {
		return nil
	}

	defer func(tx *tx) error {
		err := tx.Commit()

		if err != nil {
			log.Print(err)
		}

		return err
	}(txx)

	fn(txx)

	return nil
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
