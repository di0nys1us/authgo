package authgo

import (
	"time"
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
