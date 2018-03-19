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
	ID      int    `db:"id" json:"id,omitempty"`
	Version int    `db:"version" json:"version,omitempty"`
	Name    string `db:"name" json:"name,omitempty"`
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
