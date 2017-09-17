package repository

import (
	"database/sql"
	"time"

	"github.com/di0nys1us/valigo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	sqlFindUsers = `
		SELECT
			"id",
			"version",
			"deleted",
			"created_at",
			"created_by",
			"modified_at",
			"modified_by",
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
			"created_at",
			"created_by",
			"modified_at",
			"modified_by",
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
			"created_at",
			"created_by",
			"modified_at",
			"modified_by",
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled",
			"administrator"
		FROM "authgo"."user"
		WHERE "email" = $1;
	`
	sqlSaveUser = `
		INSERT INTO "authgo"."user" (
			"created_by",
			"modified_by",
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled",
			"administrator"
		) VALUES (
			:created_by,
			:modified_by,
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
			"modified_at" = :modified_at,
			"modified_by" = :modified_by,
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

var (
	TimeFunc = time.Now
)

type User struct {
	ID            int       `db:"id" json:"id,omitempty"`
	Version       int       `db:"version" json:"version,omitempty"`
	Deleted       bool      `db:"deleted" json:"deleted,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt,omitempty"`
	CreatedBy     string    `db:"created_by" json:"createdBy,omitempty"`
	ModifiedAt    time.Time `db:"modified_at" json:"modifiedAt,omitempty"`
	ModifiedBy    string    `db:"modified_by" json:"modifiedBy,omitempty"`
	FirstName     string    `db:"first_name" json:"firstName,omitempty"`
	LastName      string    `db:"last_name" json:"lastName,omitempty"`
	Email         string    `db:"email" json:"email,omitempty"`
	Password      string    `db:"password" json:"password,omitempty"`
	Enabled       bool      `db:"enabled" json:"enabled,omitempty"`
	Administrator bool      `db:"administrator" json:"administrator,omitempty"`
}

func (u *User) GetID() int {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) IsAdministrator() bool {
	return u.Administrator
}

func (u *User) IsInactive() bool {
	return !u.Enabled || u.Deleted
}

func (u *User) Validate() error {
	r := validation.NewValidationResult()

	if validation.IsBlank(u.CreatedBy) {
		r.AddFieldError("createdBy", "required")
	}

	if validation.IsBlank(u.ModifiedBy) {
		r.AddFieldError("modifiedBy", "required")
	}

	if validation.IsBlank(u.FirstName) {
		r.AddFieldError("firstName", "required")
	}

	if validation.IsBlank(u.LastName) {
		r.AddFieldError("lastName", "required")
	}

	if validation.IsBlank(u.Email) {
		r.AddFieldError("email", "required")
	}

	if !validation.IsEmail(u.Email) {
		r.AddFieldError("email", "email")
	}

	if validation.IsBlank(u.Password) {
		r.AddFieldError("password", "required")
	}

	if r.HasErrors() {
		return r
	}

	return nil
}

type UsersFinder interface {
	FindUsers() ([]User, error)
}

type UserFinder interface {
	FindUser(id string) (*User, error)
}

type UserByEmailFinder interface {
	FindUserByEmail(email string) (*User, error)
}

type UserSaver interface {
	SaveUser(user *User) error
}

type UserUpdater interface {
	UpdateUser(user *User) error
}

type UserDeleter interface{}

type Repository interface {
	Close() error
	UsersFinder
	UserFinder
	UserByEmailFinder
	UserSaver
	UserUpdater
	UserDeleter
}

type DefaultRepository struct {
	DB *sqlx.DB
}

func NewRepository() (Repository, error) {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		return nil, errors.Wrap(err, "authgo/repository: connection error")
	}

	return &DefaultRepository{db}, nil
}

func (r *DefaultRepository) Close() error {
	return r.DB.Close()
}

func (r *DefaultRepository) FindUsers() ([]User, error) {
	users := []User{}

	err := r.DB.Select(&users, sqlFindUsers)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/repository: error when finding users")
	}

	return users, nil
}

func (r *DefaultRepository) FindUser(id string) (*User, error) {
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

func (r *DefaultRepository) FindUserByEmail(email string) (*User, error) {
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

func (r *DefaultRepository) SaveUser(user *User) error {
	stmt, err := r.DB.PrepareNamed(sqlSaveUser)

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

func (r *DefaultRepository) UpdateUser(user *User) error {
	stmt, err := r.DB.PrepareNamed(sqlUpdateUser)

	if err != nil {
		return errors.Wrap(err, "authgo/repository: error when preparing statement")
	}

	defer stmt.Close()

	user.ModifiedAt = TimeFunc()

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
