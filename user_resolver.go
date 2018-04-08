package main

import (
	"github.com/graph-gophers/graphql-go"
)

type userResolver struct {
	repository repository
	user       *user
}

func (r *userResolver) ID() graphql.ID {
	return graphQLID(r.user.ID)
}

func (r *userResolver) Version() int32 {
	return int32(r.user.Version)
}

func (r *userResolver) FirstName() string {
	return r.user.FirstName
}

func (r *userResolver) LastName() string {
	return r.user.LastName
}

func (r *userResolver) Email() string {
	return r.user.Email
}

func (r *userResolver) Password() string {
	return r.user.Password
}

func (r *userResolver) Enabled() bool {
	return r.user.Enabled
}

func (r *userResolver) Deleted() bool {
	return r.user.Deleted
}

func (r *userResolver) Events() ([]*eventResolver, error) {
	var resolvers []*eventResolver

	for _, event := range r.user.Events {
		resolvers = append(resolvers, &eventResolver{r.repository, event})
	}

	return resolvers, nil
}

func (r *userResolver) Roles() ([]*roleResolver, error) {
	roles, err := r.repository.findUserRoles(r.user.ID.String())

	if err != nil {
		return nil, err
	}

	var resolvers []*roleResolver

	for _, role := range roles {
		resolvers = append(resolvers, &roleResolver{r.repository, role})
	}

	return resolvers, nil
}
