package main

import "github.com/pkg/errors"

type roleAuthoritiesFinder interface {
	findRoleAuthorities(roleID string) ([]*authority, error)
}

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
		return nil, errors.Wrap(err, "authgo: error when finding role authorities")
	}

	return authorities, nil
}

const (
	sqlFindRoleAuthorities = `
		select
			a.id, a.version, a.name
		from authgo.role_authority as ra
		inner join authgo.authority as a on a.id = ra.authority_id
		where ra.role_id = $1
		order by a.id;
	`
)
