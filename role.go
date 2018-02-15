package main

import (
	"database/sql"

	"github.com/pkg/errors"
)

type userRolesFinder interface {
	findUserRoles(userID string) ([]*role, error)
}

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

func (db *db) findRoleByID(id string) (*role, error) {
	r := &role{}

	err := db.Get(r, sqlFindRoleByID, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding role by id")
	}

	return r, nil
}

func (db *db) findRoleByName(name string) (*role, error) {
	r := &role{}

	err := db.Get(r, sqlFindRoleByName, name)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding role by name")
	}

	return r, nil
}

func (db *db) findAllRoles() ([]*role, error) {
	roles := []*role{}

	err := db.Select(&roles, sqlFindAllRoles)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding roles")
	}

	return roles, nil
}

func (db *db) findUserRoles(userID string) ([]*role, error) {
	roles := []*role{}

	err := db.Select(&roles, sqlFindUserRoles, userID)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding user roles")
	}

	return roles, nil
}

const (
	sqlFindRoleByID = `
		SELECT
			id,
			version,
			name
		FROM "authgo"."role"
		WHERE id = $1;
	`
	sqlFindRoleByName = `
		SELECT
			id,
			version,
			name
		FROM "authgo"."role"
		WHERE name = $1;
	`
	sqlFindAllRoles = `
		SELECT
			id,
			version,
			name
		FROM "authgo"."role"
		ORDER BY id;
	`
	sqlFindUserRoles = `
		SELECT
			r.id,
			r.version,
			r.name
		FROM "authgo"."user_role" AS ur
			INNER JOIN "authgo"."role" AS r ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.id;
	`
)
