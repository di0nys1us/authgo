package main

type repository interface {
	userRepository
	roleRepository
	authorityRepository
}

type entity struct {
	ID int `db:"id" json:"id,omitempty"`
}

type saver interface {
	save(tx *tx) error
}

type updater interface {
	update(tx *tx) error
}

type deleter interface {
	delete(tx *tx) error
}
