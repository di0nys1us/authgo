package main

type repository interface {
	userRepository
	roleRepository
	authorityRepository
	eventRepository
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
