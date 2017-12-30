package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	sqlFindUsers = `
		SELECT
			"id",
			"version",
			"deleted",
			"first_name",
			"last_name",
			"email",
			"enabled",
			"administrator"
		FROM "authgo"."user"
		ORDER BY "id";
	`
	sqlFindUser = `
		SELECT
			"id",
			"version",
			"deleted",
			"first_name",
			"last_name",
			"email",
			"enabled",
			"administrator"
		FROM "authgo"."user"
		WHERE "id" = $1;
	`
	sqlFindUserByEmail = `
		SELECT
			"id",
			"version",
			"deleted",
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled",
			"administrator"
		FROM "authgo"."user"
		WHERE "email" = $1;
	`
	sqlCreateUser = `
		INSERT INTO "authgo"."user" (
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled",
			"administrator"
		) VALUES (
			:first_name,
			:last_name,
			:email,
			:password,
			:enabled,
			:administrator
		) RETURNING "id";
	`
	sqlUpdateUser = `
		UPDATE "authgo"."user" SET
			"version" = :new_version,
			"deleted" = :deleted,
			"first_name" = :first_name,
			"last_name" = :last_name,
			"email" = :email,
			"password" = :password,
			"enabled" = :enabled,
			"administrator" = :administrator
		WHERE "id" = :id
			AND "version" = :old_version;
	`
)

// User TODO
type User struct {
	ID            int    `db:"id" json:"id,omitempty"`
	Version       int    `db:"version" json:"version,omitempty"`
	Deleted       bool   `db:"deleted" json:"deleted,omitempty"`
	FirstName     string `db:"first_name" json:"firstName,omitempty"`
	LastName      string `db:"last_name" json:"lastName,omitempty"`
	Email         string `db:"email" json:"email,omitempty"`
	Password      string `db:"password" json:"password,omitempty"`
	Enabled       bool   `db:"enabled" json:"enabled,omitempty"`
	Administrator bool   `db:"administrator" json:"administrator,omitempty"`
}

type Group struct {
	Name string `db:"name" json:"name,omitempty"`
}

type Role struct {
	Name string `db:"name" json:"name,omitempty"`
}

type GroupRole struct {
	GroupName string `db:"group_name" json:"groupName,omitempty"`
	RoleName  string `db:"role_name" json:"roleName,omitempty"`
}

// UsersFinder TODO
type UsersFinder interface {
	FindUsers() ([]User, error)
}

// UserFinder TODO
type UserFinder interface {
	FindUser(id string) (*User, error)
}

// UserByEmailFinder TODO
type UserByEmailFinder interface {
	FindUserByEmail(email string) (*User, error)
}

// UserCreator TODO
type UserCreator interface {
	CreateUser(user *User) error
}

// UserUpdater TODO
type UserUpdater interface {
	UpdateUser(user *User) error
}

// UserDeleter TODO
type UserDeleter interface{}

// Repository TODO
type Repository struct {
	DB *sqlx.DB
	UsersFinder
	UserFinder
	UserByEmailFinder
	UserCreator
	UserUpdater
	UserDeleter
}

// NewRepository TODO
func NewRepository() (*Repository, error) {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		return nil, errors.Wrap(err, "authgo/repository: connection error")
	}

	return &Repository{DB: db}, nil
}

// Close TODO
func (r *Repository) Close() error {
	return r.DB.Close()
}

// FindUsers TODO
func (r *Repository) FindUsers() ([]User, error) {
	users := []User{}

	err := r.DB.Select(&users, sqlFindUsers)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/repository: error when finding users")
	}

	return users, nil
}

// FindUser TODO
func (r *Repository) FindUser(id string) (*User, error) {
	user := &User{}

	err := r.DB.Get(user, sqlFindUser, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo/repository: error when finding user")
	}

	return user, nil
}

// FindUserByEmail TODO
func (r *Repository) FindUserByEmail(email string) (*User, error) {
	user := &User{}

	err := r.DB.Get(user, sqlFindUserByEmail, email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo/repository: error when finding user")
	}

	return user, nil
}

// CreateUser TODO
func (r *Repository) CreateUser(user *User) error {
	stmt, err := r.DB.PrepareNamed(sqlCreateUser)

	if err != nil {
		return errors.Wrap(err, "authgo/repository: error when preparing statement")
	}

	defer stmt.Close()

	var id int

	err = stmt.Get(&id, user)

	if err != nil {
		return errors.Wrap(err, "authgo/repository: error when saving user")
	}

	user.ID = id

	return nil
}

// UpdateUser TODO
func (r *Repository) UpdateUser(user *User) error {
	stmt, err := r.DB.PrepareNamed(sqlUpdateUser)

	if err != nil {
		return errors.Wrap(err, "authgo/repository: error when preparing statement")
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		struct {
			*User
			NewVersion int `db:"new_version"`
			OldVersion int `db:"old_version"`
		}{user, user.Version + 1, user.Version},
	)

	if err != nil {
		return errors.Wrap(err, "authgo/repository: error when updating user")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "authgo/repository: error when checking for rows affected")
	}

	if rowsAffected != 1 {
		return errors.New("authgo/repository: no update performed")
	}

	return nil
}
