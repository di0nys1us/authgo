package authgo

import (
	"database/sql"

	"github.com/pkg/errors"
)

type role struct {
	ID      int    `db:"id" json:"id,omitempty"`
	Version int    `db:"version" json:"version,omitempty"`
	Name    string `db:"name" json:"name,omitempty"`
}

func (r *role) save(tx *tx) error {
	return nil
}

func (r *role) update(tx *tx) error {
	return nil
}

func (r *role) delete(tx *tx) error {
	return nil
}

func findRoleByID(tx *tx, id string) (*role, error) {
	r := &role{}

	err := tx.Get(r, sqlFindRoleByID, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding role")
	}

	return r, nil
}

func findRoleByName(tx *tx, name string) (*role, error) {
	r := &role{}

	err := tx.Get(r, sqlFindRoleByName, name)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding role")
	}

	return r, nil
}

const (
	sqlFindRoleByID = `
		SELECT
			r.id,
			r.version,
			r.name
		FROM "authgo"."role" AS r
		WHERE r.id = $1;
	`
	sqlFindRoleByName = `
		SELECT
			r.id,
			r.version,
			r.name
		FROM "authgo"."role" AS r
		WHERE r.name = $1;
	`
	sqlFindAllRoles = `
		SELECT
			r.id,
			r.version,
			r.name
		FROM "authgo"."role" AS r
		ORDER BY r.id;
	`
	sqlFindUserRoles = `
		SELECT
			r.id,
			r.version,
			r.name,
		FROM "authgo"."user_role" AS ur
			INNER JOIN "authgo"."role" AS r ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.id;
	`
)
