package repository

import (
	"github.com/jmoiron/sqlx"
)

const (
	sqlFindRoleByID = `
		SELECT
			r.id,
			r.name
		FROM "authgo"."role" AS r
		WHERE r.id = $1;
	`
	sqlFindRoleByName = `
		SELECT
			r.id,
			r.name
		FROM "authgo"."role" AS r
		WHERE r.name = $1;
	`
	sqlFindAllRoles = `
		SELECT
			r.id,
			r.name
		FROM "authgo"."role" AS r
		ORDER BY r.id;
	`
	sqlFindUserRoles = `
		SELECT
			r.id,
			r.name,
		FROM "authgo"."user_role" AS ur
			INNER JOIN "authgo"."role" AS r ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.id;
	`
)

// Role TODO
type Role struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

// RoleRepository TODO
type RoleRepository interface {
	FindByID(id int) (*Role, error)
	FindByName(name int) (*Role, error)
	FindAll() ([]*Role, error)
	FindByUserID(userID int) ([]*Role, error)
}

type roleRepository struct {
	RoleRepository
	db *sqlx.DB
}
