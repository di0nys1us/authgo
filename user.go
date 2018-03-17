package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// INTERFACES

type userByIDFinder interface {
	findUserByID(id string) (*user, error)
}

type userByEmailFinder interface {
	findUserByEmail(email string) (*user, error)
}

type allUsersFinder interface {
	findAllUsers() ([]*user, error)
}

type userSaver interface {
	saveUser(*user) error
}

type userRepository interface {
	userByIDFinder
	userByEmailFinder
	allUsersFinder
	userSaver
}

// STRUCTS

type user struct {
	ID        int    `db:"id" json:"id,omitempty"`
	Version   int    `db:"version" json:"version,omitempty"`
	FirstName string `db:"first_name" json:"firstName,omitempty"`
	LastName  string `db:"last_name" json:"lastName,omitempty"`
	Email     string `db:"email" json:"email,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Enabled   bool   `db:"enabled" json:"enabled,omitempty"`
	Deleted   bool   `db:"deleted" json:"deleted,omitempty"`
}

func (u *user) save(tx *tx) error {
	id, err := tx.save(u, sqlSaveUser)

	if err != nil {
		return errors.WithStack(err)
	}

	u.ID = id

	return nil
}

func (u *user) update(tx *tx) error {
	stmt, err := tx.PrepareNamed(sqlUpdateUser)

	if err != nil {
		return errors.WithStack(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(
		struct {
			*user
			NewVersion int `db:"new_version"`
			OldVersion int `db:"old_version"`
		}{u, u.Version + 1, u.Version},
	)

	if err != nil {
		return errors.WithStack(err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return errors.WithStack(err)
	}

	if rowsAffected != 1 {
		return errors.New("authgo: no update performed")
	}

	return nil
}

func (u *user) delete(tx *tx) error {
	return nil
}

func (db *db) findUserByID(id string) (*user, error) {
	u := &user{}

	err := db.Get(u, sqlFindUserByID, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}

func (db *db) findUserByEmail(email string) (*user, error) {
	u := &user{}

	err := db.Get(u, sqlFindUserByEmail, email)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u, nil
}

func (db *db) findAllUsers() ([]*user, error) {
	users := []*user{}

	err := db.Select(&users, sqlFindAllUsers)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return users, nil
}

func (db *db) saveUser(user *user) error {
	return db.run(func(tx *tx) error {
		err := user.save(tx)

		if err != nil {
			return errors.WithStack(err)
		}

		event := &event{
			CreatedBy:   1,
			CreatedAt:   time.Now(),
			Type:        eventTypeUserCreated,
			Description: fmt.Sprintf("User %q created.", user.Email),
		}

		err = event.save(tx)

		if err != nil {
			return errors.WithStack(err)
		}

		userEvent := &userEvent{user.ID, event.ID}

		err = userEvent.save(tx)

		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

const (
	sqlSaveUser = `
		INSERT INTO "authgo"."user" (
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled",
			"deleted"
		) VALUES (
			:first_name,
			:last_name,
			:email,
			:password,
			:enabled,
			:deleted
		) RETURNING "id";
	`
	sqlUpdateUser = `
		UPDATE "authgo"."user" SET
			"version" = :new_version,
			"first_name" = :first_name,
			"last_name" = :last_name,
			"email" = :email,
			"password" = :password,
			"enabled" = :enabled,
			"deleted" = :deleted
		WHERE "id" = :id
			AND "version" = :old_version;
	`
	sqlDeleteUser   = ``
	sqlFindUserByID = `
		SELECT
			"id",
			"version",
			"first_name",
			"last_name",
			"email",
			"enabled",
			"deleted"
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
			"enabled",
			"deleted"
		FROM "authgo"."user"
		WHERE "email" = $1;
	`
	sqlFindAllUsers = `
		SELECT
			"id",
			"version",
			"first_name",
			"last_name",
			"email",
			"enabled",
			"deleted"
		FROM "authgo"."user"
		ORDER BY "id";
	`
)
