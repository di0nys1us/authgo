package main

import (
	"time"

	"github.com/graph-gophers/graphql-go"
)

type eventResolver struct {
	repository repository
	event      *event
}

func (r *eventResolver) ID() graphql.ID {
	return graphQLID(r.event.ID)
}

func (r *eventResolver) CreatedBy() (*userResolver, error) {
	user, err := r.repository.findUserByID(r.event.CreatedBy)

	if err != nil {
		return nil, err
	}

	return &userResolver{r.repository, user}, nil
}

func (r *eventResolver) CreatedAt() string {
	return r.event.CreatedAt.Format(time.RFC3339)
}

func (r *eventResolver) Type() string {
	return r.event.Type
}

func (r *eventResolver) Description() string {
	return r.event.Description
}
