package authgo

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

type event struct {
	ID          int       `db:"id" json:"id,omitempty"`
	CreatedBy   int       `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	Type        int       `db:"type" json:"type,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
}

func (evt *event) save(tx *tx) error {
	return nil
}

func (evt *event) update(tx *tx) error {
	return nil
}

func (evt *event) delete(tx *tx) error {
	return nil
}

func findEventByID(tx *tx, id string) (*event, error) {
	evt := &event{}

	err := tx.Get(evt, sqlFindEventByID, id)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding event")
	}

	return evt, nil
}

const (
	sqlFindEventByID = `
		SELECT id, created_by, created_at, type, description
			FROM authgo.event WHERE id = $1;
	`
)
