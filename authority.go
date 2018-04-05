package main

import (
	"github.com/pkg/errors"
)

// INTERFACES

type roleAuthoritiesFinder interface {
	findRoleAuthorities(roleID string) ([]*authority, error)
}

type authorityRepository interface {
	roleAuthoritiesFinder
}

// STRUCTS

type authority struct {
	ID      string   `db:"id"`
	Version int      `db:"version"`
	Events  []*event `db:"events"`
	Name    string   `db:"name"`
}

func (a *authority) save(tx *tx) error {
	return nil
}

func (a *authority) update(tx *tx) error {
	return nil
}

func (a *authority) delete(tx *tx) error {
	return nil
}

func (db *db) findRoleAuthorities(roleID string) ([]*authority, error) {
	authorities := []*authority{}

	err := db.Select(&authorities, sqlFindRoleAuthorities, roleID)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return authorities, nil
}

const (
	sqlFindRoleAuthorities = `
		select
			"authority"."id",
			"authority"."version",
			"authority"."name"
		from "authgo"."authority"
		inner join "authgo"."role_authority" on "role_authority"."authority_id" = "authority"."id"
		where "role_authority"."role_id" = $1
		order by "authority"."id";
	`
)
