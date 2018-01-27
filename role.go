package authgo

type role struct {
	ID      int    `db:"id" json:"id,omitempty"`
	Version int    `db:"version" json:"version,omitempty"`
	Name    string `db:"name" json:"name,omitempty"`
}

func (r *role) create(tx *tx) error {
	return nil
}

func (r *role) update(tx *tx) error {
	return nil
}

func (r *role) delete(tx *tx) error {
	return nil
}
