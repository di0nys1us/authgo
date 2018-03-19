package main

import (
	graphql "github.com/neelance/graphql-go"
)

type rootQuery struct {
	repository repository
}

func (r *rootQuery) Users() ([]*userResolver, error) {
	users, err := r.repository.findAllUsers()

	if err != nil {
		return nil, err
	}

	var resolvers []*userResolver

	for _, user := range users {
		resolvers = append(resolvers, &userResolver{r.repository, user})
	}

	return resolvers, nil
}

func (r *rootQuery) User(args struct {
	ID    *graphql.ID
	Email *string
}) (*userResolver, error) {
	var user *user
	var err error

	if args.ID != nil {
		user, err = r.repository.findUserByID(string(*args.ID))
	}

	if args.Email != nil {
		user, err = r.repository.findUserByEmail(*args.Email)
	}

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &userResolver{r.repository, user}, nil
}

func (r *rootQuery) Events(args struct {
	UserID *graphql.ID
}) ([]*eventResolver, error) {
	var events []*event
	var err error

	if args.UserID != nil {
		events, err = r.repository.findEventsByUserID(string(*args.UserID))
	} else {
		events, err = r.repository.findAllEvents()
	}

	if err != nil {
		return nil, err
	}

	var resolvers []*eventResolver

	for _, event := range events {
		resolvers = append(resolvers, &eventResolver{r.repository, event})
	}

	return resolvers, nil
}

func (r *rootQuery) Event(args struct {
	ID graphql.ID
}) (*eventResolver, error) {
	event, err := r.repository.findEventByID(string(args.ID))

	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, nil
	}

	return &eventResolver{r.repository, event}, nil
}
