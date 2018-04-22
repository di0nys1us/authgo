package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/di0nys1us/authgo/security"
	"github.com/pkg/errors"
)

// INTERFACES

type allUsersFinder interface {
	findAllUsers() ([]*user, error)
}

type userByIDFinder interface {
	findUserByID(id string) (*user, error)
}

type userByEmailFinder interface {
	findUserByEmail(email string) (*user, error)
}

type userSaver interface {
	saveUser(ctx context.Context, user *user) error
}

type userRepository interface {
	allUsersFinder
	userByIDFinder
	userByEmailFinder
	userSaver
}

// STRUCTS

type user struct {
	ID        string `db:"id" json:"id,omitempty"`
	Version   int    `db:"version" json:"version,omitempty"`
	FirstName string `db:"first_name" json:"firstName,omitempty"`
	LastName  string `db:"last_name" json:"lastName,omitempty"`
	Email     string `db:"email" json:"email,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Enabled   bool   `db:"enabled" json:"enabled,omitempty"`
	Deleted   bool   `db:"deleted" json:"deleted,omitempty"`
	Events    events `db:"events" json:"events,omitempty"`
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

func (db *db) findAllUsers() ([]*user, error) {
	users := []*user{}

	err := db.Select(&users, sqlFindAllUsers)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return users, nil
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

func (db *db) saveUser(ctx context.Context, user *user) error {
	return db.commit(func(tx *tx) error {
		eventID, err := db.generateUUID()

		if err != nil {
			return errors.WithStack(err)
		}

		userID, err := security.UserIDFromContext(ctx)

		if err != nil {
			return errors.WithStack(err)
		}

		event := &event{
			ID:          eventID,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			Type:        eventTypeUserCreated,
			Description: fmt.Sprintf("User %q created.", user.Email),
		}

		user.Events = append(user.Events, event)

		err = user.save(tx)

		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

const (
	sqlSaveUser = `
		insert into "authgo"."user" (
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled",
			"deleted",
			"events"
		) values (
			:first_name,
			:last_name,
			:email,
			:password,
			:enabled,
			:deleted,
			:events
		) returning "user"."id";
	`
	sqlUpdateUser = `
		update "authgo"."user" set
			"version" = :new_version,
			"first_name" = :first_name,
			"last_name" = :last_name,
			"email" = :email,
			"password" = :password,
			"enabled" = :enabled,
			"deleted" = :deleted
		where "user"."id" = :id
			and "user"."version" = :old_version;
	`
	sqlDeleteUser   = ``
	sqlFindUserByID = `
		select
			"user"."id",
			"user"."version",
			"user"."first_name",
			"user"."last_name",
			"user"."email",
			"user"."enabled",
			"user"."deleted"
		from "authgo"."user"
		where "user"."id" = $1;
	`
	sqlFindUserByEmail = `
		select
			"user"."id",
			"user"."version",
			"user"."first_name",
			"user"."last_name",
			"user"."email",
			"user"."password",
			"user"."enabled",
			"user"."deleted"
		from "authgo"."user"
		where "user"."email" = $1;
	`
	sqlFindAllUsers = `
		select
			"user"."id",
			"user"."version",
			"user"."first_name",
			"user"."last_name",
			"user"."email",
			"user"."enabled",
			"user"."deleted",
			"user"."events"
		from "authgo"."user"
		order by "user"."id";
	`
	sqlFindRoleUsers = `
		select
			"user"."id",
			"user"."version",
			"user"."first_name",
			"user"."last_name",
			"user"."email",
			"user"."enabled",
			"user"."deleted",
			"user"."events"
		from "authgo"."user"
			inner join "authgo"."user_role" on "user_role"."user_id" = "user"."id"
		where "user_role"."role_id" = $1
		order by "user"."id";
	`
)
