package authgo

const (
	eventTypeCreateUser = "CREATE_USER"
	eventTypeUpdateUser = "UPDATE_USER"
)

type eventType struct {
	ID   int    `db:"id" json:"id,omitempty"`
	Name string `db:"name" json:"name,omitempty"`
}
