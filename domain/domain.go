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

// User TODO
type User struct {
	ID        int    `db:"id" json:"id,omitempty"`
	Version   int    `db:"version" json:"version,omitempty"`
	FirstName string `db:"first_name" json:"firstName,omitempty"`
	LastName  string `db:"last_name" json:"lastName,omitempty"`
	Email     string `db:"email" json:"email,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Enabled   bool   `db:"enabled" json:"enabled,omitempty"`
}

// Role TODO
type Role struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}

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
