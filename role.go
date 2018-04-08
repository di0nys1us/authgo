package main

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type userRolesFinder interface {
	findUserRoles(userID string) ([]*role, error)
}

type roleRepository interface {
	userRolesFinder
}

type role struct {
	ID      uuid.UUID `db:"id" json:"id,omitempty"`
	Version int       `db:"version" json:"version,omitempty"`
	Name    string    `db:"name" json:"name,omitempty"`
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
		select
			"role"."id",
			"role"."version",
			"role"."name"
		from "authgo"."role"
		where "role"."id" = $1;
	`
	sqlFindRoleByName = `
		select
			"role"."id",
			"role"."version",
			"role"."name"
		from "authgo"."role"
		where "role"."name" = $1;
	`
	sqlFindAllRoles = `
		select
			"role"."id",
			"role"."version",
			"role"."name"
		from "authgo"."role"
		order by "role"."id";
	`
	sqlFindUserRoles = `
		select
			"role"."id",
			"role"."version",
			"role"."name"
		from "authgo"."role"
			inner join "authgo"."user_role" on "user_role"."role_id" = "role"."id"
		where "user_role"."user_id" = $1
		order by "role"."id";
	`
)
