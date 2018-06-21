package model

// Role TODO
type Role struct {
	ID      string `db:"id"`
	Version int    `db:"version"`
	Name    string `db:"name"`
	Events  Events `db:"events" json:"events,omitempty"`
}
