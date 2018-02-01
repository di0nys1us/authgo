package authgo

type authority struct {
	ID      int    `db:"id" json:"id,omitempty"`
	Version int    `db:"version" json:"version,omitempty"`
	Name    string `db:"name" json:"name,omitempty"`
}

func (a *authority) save(tx *tx) error {
	return nil
}

func (a *authority) update(tx *tx) error {
	return nil
}

func (a *authority) delete(tx *tx) error {
	return nil
}
