package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (db *db) saveAndCommit(saver ...saver) error {
	tx, err := db.save(saver...)

	if err != nil {
		return err
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (db *db) save(saver ...saver) (*tx, error) {
	tx, err := db.begin()

	if err != nil {
		return nil, err
	}

	for _, s := range saver {
		err = s.save(tx)

		if err != nil {
			return nil, err
		}
	}

	return tx, nil
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

func (tx *tx) saveEntity(arg interface{}, query string) (*entity, error) {
	stmt, err := tx.PrepareNamed(query)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when preparing statement")
	}

	defer stmt.Close()

	var id int

	err = stmt.Get(&id, arg)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when saving")
	}

	return &entity{id}, nil
}

func (tx *tx) save(arg interface{}, query string) error {
	stmt, err := tx.PrepareNamed(query)

	if err != nil {
		return errors.Wrap(err, "authgo: error when preparing statement")
	}

	defer stmt.Close()

	_, err = stmt.Exec(arg)

	if err != nil {
		return errors.Wrap(err, "authgo: error when saving")
	}

	return nil
}
