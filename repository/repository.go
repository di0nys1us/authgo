package repository

import (
	"database/sql"

	"github.com/di0nys1us/authgo/domain"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	sqlFindUsers = `
		SELECT
			"id",
			"version",
			"first_name",
			"last_name",
			"email",
			"enabled"
		FROM "authgo"."user"
		ORDER BY "id";
	`
	sqlFindUser = `
		SELECT
			"id",
			"version",
			"first_name",
			"last_name",
			"email",
			"enabled"
		FROM "authgo"."user"
		WHERE "id" = $1;
	`
	sqlFindUserByEmail = `
		SELECT
			"id",
			"version",
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled"
		FROM "authgo"."user"
		WHERE "email" = $1;
	`
	sqlCreateUser = `
		INSERT INTO "authgo"."user" (
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled"
		) VALUES (
			:first_name,
			:last_name,
			:email,
			:password,
			:enabled
		) RETURNING "id";
	`
	sqlUpdateUser = `
		UPDATE "authgo"."user" SET
			"version" = :new_version,
			"first_name" = :first_name,
			"last_name" = :last_name,
			"email" = :email,
			"password" = :password,
			"enabled" = :enabled
		WHERE "id" = :id
			AND "version" = :old_version;
	`
)

// UsersFinder TODO
type UsersFinder interface {
	FindUsers() ([]domain.User, error)
}

// UserFinder TODO
type UserFinder interface {
	FindUser(id string) (*domain.User, error)
}

// UserByEmailFinder TODO
type UserByEmailFinder interface {
	FindUserByEmail(email string) (*domain.User, error)
}

// UserCreator TODO
type UserCreator interface {
	CreateUser(user *domain.User) error
}

// UserUpdater TODO
type UserUpdater interface {
	UpdateUser(user *domain.User) error
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

func CreateDatabase() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
}

// NewRepository TODO
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{DB: db}
}

// FindUsers TODO
func (r *Repository) FindUsers() ([]domain.User, error) {
	users := []domain.User{}

	err := r.DB.Select(&users, sqlFindUsers)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/repository: error when finding users")
	}

	return users, nil
}

// FindUser TODO
func (r *Repository) FindUser(id string) (*domain.User, error) {
	user := &domain.User{}

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
func (r *Repository) FindUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

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
func (r *Repository) CreateUser(user *domain.User) error {
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
func (r *Repository) UpdateUser(user *domain.User) error {
	stmt, err := r.DB.PrepareNamed(sqlUpdateUser)

	if err != nil {
		return errors.Wrap(err, "authgo/repository: error when preparing statement")
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		struct {
			*domain.User
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
