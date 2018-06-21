package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/di0nys1us/authgo/security"
	uuid "github.com/satori/go.uuid"
)

// Event TODO
type Event struct {
	ID          string    `json:"id,omitempty"`
	CreatedBy   string    `json:"createdBy,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	Type        string    `json:"type,omitempty"`
	Description string    `json:"description,omitempty"`
}

// Events TODO
type Events []*Event

// Scan TODO
func (e *Events) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, e)
	}

	return nil
}

// Value TODO
func (e *Events) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// NewEvent TODO
func NewEvent(ctx context.Context, eventType, description string) (*Event, error) {
	eventID, err := uuid.NewV1()

	if err != nil {
		return nil, err
	}

	return &Event{
		ID:          eventID.String(),
		CreatedBy:   security.UserEmailFromContext(ctx),
		CreatedAt:   time.Now(),
		Type:        eventType,
		Description: description,
	}, nil
}
