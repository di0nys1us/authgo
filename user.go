package authgo

import (
	"database/sql"

	"github.com/pkg/errors"
)

type user struct {
	ID        int    `db:"id" json:"id,omitempty"`
	Version   int    `db:"version" json:"version,omitempty"`
	FirstName string `db:"first_name" json:"firstName,omitempty"`
	LastName  string `db:"last_name" json:"lastName,omitempty"`
	Email     string `db:"email" json:"email,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Enabled   bool   `db:"enabled" json:"enabled,omitempty"`
}

func (u *user) save(tx *tx) error {
	return nil
}

func (u *user) update(tx *tx) error {
	return nil
}

func (u *user) delete(tx *tx) error {
	return nil
}

func findUserByID(tx *tx, id int) (*user, error) {
	u := &user{}

	err := tx.Get(u, sqlFindUser, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding user")
	}

	return u, nil
}

func findUserByEmail(tx *tx, email string) (*user, error) {
	u := &user{}

	err := tx.Get(u, sqlFindUser, email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding user")
	}

	return u, nil
}

func findAllUsers(tx *tx) ([]*user, error) {
	return nil, nil
}

const (
	sqlFindUser = `
		SELECT
			"id",
			"version",
			"first_name",
			"last_name",
			"email",
			"enabled"
		FROM "authgo"."user"
		WHERE "id" = $1;
	`
)
