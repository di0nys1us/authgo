package model

// Authority TODO
type Authority struct {
	ID      string `db:"id"`
	Version int    `db:"version"`
	Name    string `db:"name"`
	Events  Events `db:"events" json:"events,omitempty"`
}
