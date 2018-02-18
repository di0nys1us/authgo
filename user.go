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
	*entity
	Version   int    `db:"version" json:"version,omitempty"`
	FirstName string `db:"first_name" json:"firstName,omitempty"`
	LastName  string `db:"last_name" json:"lastName,omitempty"`
	Email     string `db:"email" json:"email,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Enabled   bool   `db:"enabled" json:"enabled,omitempty"`
	Deleted   bool   `db:"deleted" json:"deleted,omitempty"`
}

func (u *user) save(tx *tx) error {
	entity, err := tx.saveEntity(u, sqlSaveUser)

	if err != nil {
		return errors.Wrap(err, "authgo: error when saving user")
	}

	u.entity = entity

	return nil
}

func (u *user) update(tx *tx) error {
	stmt, err := tx.PrepareNamed(sqlUpdateUser)

	if err != nil {
		return errors.Wrap(err, "authgo: error when preparing statement")
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
		return errors.Wrap(err, "authgo: error when updating user")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return errors.Wrap(err, "authgo: error when checking for rows affected")
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
		return nil, errors.Wrap(err, "authgo: error when finding user by id")
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
		return nil, errors.Wrap(err, "authgo: error when finding user by email")
	}

	return u, nil
}

func (db *db) findAllUsers() ([]*user, error) {
	users := []*user{}

	err := db.Select(&users, sqlFindAllUsers)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding users")
	}

	return users, nil
}

func (db *db) saveUser(user *user) error {
	event := &event{
		CreatedBy:   1,
		CreatedAt:   time.Now(),
		Type:        eventTypeUserCreated,
		Description: fmt.Sprintf("User %q created.", user.Email),
	}

	tx, err := db.save(user, event)

	userEvent := &userEvent{user.ID, event.ID}

	userEvent.save(tx)

	tx.Commit()

	if err != nil {
		return err
	}

	return nil
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
