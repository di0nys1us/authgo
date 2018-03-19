package main

import (
	"strconv"

	graphql "github.com/neelance/graphql-go"
)

type eventResolver struct {
	repository repository
	event      *event
}

func (r *eventResolver) ID() graphql.ID {
	return intToID(r.event.ID)
}

func (r *eventResolver) CreatedBy() (*userResolver, error) {
	user, err := r.repository.findUserByID(strconv.Itoa(r.event.CreatedBy))

	if err != nil {
		return nil, err
	}

	return &userResolver{r.repository, user}, nil
}

func (r *eventResolver) CreatedAt() string {
	return r.event.CreatedAt.String()
}

func (r *eventResolver) Type() (string, error) {
	return r.event.Type, nil
}

func (r *eventResolver) Description() string {
	return r.event.Description
}
