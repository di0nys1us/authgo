package main

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// INTERFACES

type allEventsFinder interface {
	findAllEvents() ([]*event, error)
}

type eventByIDFinder interface {
	findEventByID(id string) (*event, error)
}

type eventsByUserIDFinder interface {
	findEventsByUserID(userID string) ([]*event, error)
}

type eventRepository interface {
	allEventsFinder
	eventByIDFinder
	eventsByUserIDFinder
}

// STRUCTS

type event struct {
	ID          int       `db:"id" json:"id,omitempty"`
	CreatedBy   int       `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	Type        string    `db:"type" json:"type,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
}

func (e *event) save(tx *tx) error {
	id, err := tx.save(e, sqlSaveEvent)

	if err != nil {
		return errors.WithStack(err)
	}

	e.ID = id

	return nil
}

type userEvent struct {
	UserID  int `db:"user_id"`
	EventID int `db:"event_id"`
}

func (e *userEvent) save(tx *tx) error {
	_, err := tx.save(e, sqlSaveUserEvent)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (db *db) findAllEvents() ([]*event, error) {
	events := []*event{}

	err := db.Select(&events, sqlFindAllEvents)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return events, nil
}

func (db *db) findEventByID(id string) (*event, error) {
	event := &event{}

	err := db.Get(event, sqlFindEventByID, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return event, nil
}

func (db *db) findEventsByUserID(userID string) ([]*event, error) {
	events := []*event{}

	err := db.Select(&events, sqlFindEventsByUserID, userID)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return events, nil
}

func (db *db) findEventsByRoleID(roleID string) ([]*event, error) {
	return nil, nil
}

func (db *db) findEventsByAuthorityID(authorityID string) ([]*event, error) {
	return nil, nil
}

const (
	sqlSaveEvent = `
		insert into authgo.event (created_by, created_at, type, description)
		values (:created_by, :created_at, :type, :description)
		returning id;
	`
	sqlSaveUserEvent = `
		insert into authgo.user_event (user_id, event_id)
		values (:user_id, :event_id)
		returning -1 as id;
	`
	sqlFindAllEvents = `
		select
			"event"."id",
			"event"."created_by",
			"event"."created_at",
			"event"."type",
			"event"."description"
		from "authgo"."event"
		order by "event"."created_at" desc, "event"."id" desc;
	`
	sqlFindEventByID = `
		select
			"event"."id",
			"event"."created_by",
			"event"."created_at",
			"event"."type",
			"event"."description"
		from "authgo"."event"
		where "event"."id" = $1;
	`
	sqlFindEventsByUserID = `
		select
			"event"."id",
			"event"."created_by",
			"event"."created_at",
			"event"."type",
			"event"."description"
		from "authgo"."event"
			inner join "authgo"."user_event" on "user_event"."event_id" = "event"."id"
		where "user_event"."user_id" = $1
		order by "event"."created_at" desc, "event"."id" desc;
	`
)

const (
	eventTypeUserCreated = "USER_CREATED"
	eventTypeUserUpdated = "USER_UPDATED"
)

type eventType struct {
	Name        string `db:"name" json:"name,omitempty"`
	Description string `db:"description" json:"description,omitempty"`
}
