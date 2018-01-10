package models

const (
	EventTypeCreateUser = "CREATE_USER"
	EventTypeUpdateUser = "UPDATE_USER"
)

type UserEvent struct {
	ID     int
	UserID int
}
