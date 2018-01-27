package domain

import (
	"time"
)

const (
	// EventTypeCreateUser TODO
	EventTypeCreateUser = "CREATE_USER"
	// EventTypeUpdateUser TODO
	EventTypeUpdateUser = "UPDATE_USER"
)

// Authority TODO
type Authority struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

// EventType TODO
type EventType struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

// UserEvent TODO
type UserEvent struct {
	ID          int       `db:"id" json:"id,omitempty"`
	UserID      int       `db:"user_id" json:"user_id,omitempty"`
	CreatedBy   int       `db:"created_by" json:"created_by,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at,omitempty"`
	Type        int       `db:"type" json:"type,omitempty"`
	Description string    `db:"description" json:"description,omitempty"`
}
