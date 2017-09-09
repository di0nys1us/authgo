package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/di0nys1us/valigo"
	"github.com/jmoiron/sqlx"
)

var (
	TimeFunc             = time.Now
	ErrNoUpdatePerformed = errors.New("authgo/repository: no update")
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

func (u *User) Validate() (bool, error) {
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

	return !r.HasFieldErrors(), r
}

type Users []*User

type UsersFinder interface {
	FindUsers() (*Users, error)
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
		return nil, err
	}

	return &DefaultRepository{db}, nil
}

func (r *DefaultRepository) FindUsers() (*Users, error) {
	const q = `
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

	users := &Users{}

	err := r.DB.Select(users, q)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *DefaultRepository) FindUser(id string) (*User, error) {
	const q = `
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

	user := &User{}

	err := r.DB.Get(user, q, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *DefaultRepository) FindUserByEmail(email string) (*User, error) {
	const q = `
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

	user := &User{}

	err := r.DB.Get(user, q, email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *DefaultRepository) SaveUser(user *User) error {
	const q = `
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

	stmt, err := r.DB.PrepareNamed(q)

	if err != nil {
		return err
	}

	defer stmt.Close()

	var id int

	err = stmt.Get(&id, user)

	if err != nil {
		return err
	}

	user.ID = id

	return nil
}

func (r *DefaultRepository) UpdateUser(user *User) error {
	const q = `
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
		WHERE "id" = :id AND "version" = :old_version;
	`

	stmt, err := r.DB.PrepareNamed(q)

	if err != nil {
		return err
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
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return ErrNoUpdatePerformed
	}

	return nil
}
