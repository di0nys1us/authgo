package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/di0nys1us/authgo/security"
	"github.com/di0nys1us/authgo/sqlgo"
)

// User TODO
type User struct {
	ID             string `json:"id,omitempty"`
	Version        int    `json:"version"`
	Events         Events `json:"events,omitempty"`
	Deleted        bool   `json:"deleted"`
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"-"`
	HashedPassword string `json:"-"`
	Enabled        bool   `json:"enabled"`
}

// Find a user
func (u *User) Find(q sqlgo.QueryerRow) error {
	var row *sql.Row

	switch {
	case u.ID != "":
		const query = ` 
			select
				"user"."id",
				"user"."version",
				"user"."events",
				"user"."deleted",
				"user"."first_name",
				"user"."last_name",
				"user"."email",
				"user"."password",
				"user"."enabled"
			from "authgo"."user"
			where "user"."id" = $1;
		`

		row = q.QueryRow(query, u.ID)
	case u.Email != "":
		const query = `
			select
				"user"."id",
				"user"."version",
				"user"."events",
				"user"."deleted",
				"user"."first_name",
				"user"."last_name",
				"user"."email",
				"user"."password",
				"user"."enabled"
			from "authgo"."user"
			where "user"."email" = $1;
		`

		row = q.QueryRow(query, u.Email)
	default:
		return errors.New("authgo: missing both id and email")
	}

	err := row.Scan(
		&u.ID,
		&u.Version,
		&u.Events,
		&u.Deleted,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.Enabled,
	)

	if err != nil {
		return err
	}

	return nil
}

// Save the user
func (u *User) Save(ctx context.Context, tx *sqlgo.Tx) error {
	const query = `
		insert into "authgo"."user" (
			"events",
			"deleted",
			"first_name",
			"last_name",
			"email",
			"password",
			"enabled"
		) values (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7
		) returning "user"."id";
	`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer sqlgo.Close(stmt)

	event, err := NewEvent(ctx, "USER_CREATED", fmt.Sprintf("User %q created.", u.Email))

	if err != nil {
		return err
	}

	events := append(u.Events, event)

	hashedPassword, err := security.GenerateHashedPassword(u.Password)

	if err != nil {
		return err
	}

	row := stmt.QueryRow(
		&events,
		u.Deleted,
		u.FirstName,
		u.LastName,
		u.Email,
		hashedPassword,
		u.Enabled,
	)

	err = row.Scan(&u.ID)

	if err != nil {
		return err
	}

	u.Events = events
	u.HashedPassword = hashedPassword

	return nil
}

// Delete the user
func (u *User) Delete(tx *sqlgo.Tx) error {
	const query = `
		update "authgo"."user"
		set "deleted" = true 
		where "user"."id" = $1;
	`

	stmt, err := tx.Prepare(query)

	if err != nil {
		return err
	}

	defer sqlgo.Close(stmt)

	_, err = stmt.Exec(u.ID)

	if err != nil {
		return err
	}

	u.Deleted = true

	return nil
}

// Users TODO
type Users []*User

// FindAll TODO
func (u *Users) FindAll(tx *sqlgo.Tx) error {
	const query = `
		select
			"user"."id",
			"user"."version",
			"user"."events",
			"user"."deleted",
			"user"."first_name",
			"user"."last_name",
			"user"."email",
			"user"."password",
			"user"."enabled"
		from "authgo"."user";
	`

	rows, err := tx.Query(query)

	if err != nil {
		return err
	}

	defer sqlgo.Close(rows)

	for rows.Next() {
		user := &User{}

		err = rows.Scan(
			&user.ID,
			&user.Version,
			&user.Events,
			&user.Deleted,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.Enabled,
		)

		if err != nil {
			return err
		}

		*u = append(*u, user)
	}

	err = rows.Err()

	if err != nil {
		return err
	}

	return nil
}
