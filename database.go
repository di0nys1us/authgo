package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func newDB() (*db, error) {
	wrapped, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &db{wrapped}, nil
}

type db struct {
	*sqlx.DB
}

func (db *db) begin() (*tx, error) {
	wrapped, err := db.DB.Beginx()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &tx{wrapped}, nil
}

func (db *db) run(fn func(tx *tx) error) error {
	tx, err := db.begin()

	if err != nil {
		return errors.WithStack(err)
	}

	err = fn(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Commit()

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type tx struct {
	*sqlx.Tx
}

func (tx *tx) save(arg interface{}, query string) (int, error) {
	stmt, err := tx.PrepareNamed(query)

	if err != nil {
		return -1, errors.WithStack(err)
	}

	defer stmt.Close()

	var id int

	err = stmt.Get(&id, arg)

	if err != nil {
		return -1, errors.WithStack(err)
	}

	return id, nil
}
