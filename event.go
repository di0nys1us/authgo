package main

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// INTERFACES

type eventRepository interface {
}

// STRUCTS

type event struct {
	*entity
	CreatedBy   int       `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	Type        string    `db:"type" json:"type,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
}

func (e *event) save(tx *tx) error {
	entity, err := tx.saveEntity(e, sqlSaveEvent)

	if err != nil {
		return errors.Wrap(err, "authgo: error when saving event")
	}

	e.entity = entity

	return nil
}

func (e *event) update(tx *tx) error {
	return nil
}

func (e *event) delete(tx *tx) error {
	return nil
}

type userEvent struct {
	UserID  int `db:"user_id"`
	EventID int `db:"event_id"`
}

func (e *userEvent) save(tx *tx) error {
	err := tx.save(e, sqlSaveUserEvent)

	if err != nil {
		return err
	}

	return nil
}

func (db *db) findEventByID(id string) (*event, error) {
	e := &event{}

	err := db.Get(e, sqlFindEventByID, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding event by id")
	}

	return e, nil
}

func (db *db) findEventsByUserID(userID string) ([]*event, error) {
	return nil, nil
}

const (
	sqlSaveEvent = `
		insert into "authgo"."event" ("created_by", "created_at", "type", "description")
		values (:created_by, :created_at, :type, :description)
		returning "id";
	`
	sqlSaveUserEvent = `
		insert into "authgo"."event" ("user_id", "event_id")
		values (:user_id, :event_id)
	`
	sqlFindEventByID = `
		SELECT
			id, created_by, created_at, type, description
		FROM authgo.event WHERE id = $1;
	`
	sqlFindAllEvents = `
		SELECT
			id, created_by, created_at, type, description
		FROM authgo.event
		ORDER BY id;
	`
	sqlFindEventsByUserID = `
		select
			e.id,
			e.created_by,
			e.created_at,
			e.type,
			e.description
		from authgo."event" as e
		inner join authgo.user_event as ue on e.id = ue.user_id
		where ue.user_id = 1
		order by e.created_at desc, e.id desc;
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
