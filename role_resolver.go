package main

import (
	"github.com/graph-gophers/graphql-go"
)

type roleResolver struct {
	repository repository
	role       *role
}

func (r *roleResolver) ID() graphql.ID {
	return graphQLID(r.role.ID)
}

func (r *roleResolver) Version() int32 {
	return int32(r.role.Version)
}

func (r *roleResolver) Name() string {
	return r.role.Name
}

func (r *roleResolver) Events() ([]*eventResolver, error) {
	return nil, nil
}

func (r *roleResolver) Authorities() ([]*authorityResolver, error) {
	authorities, err := r.repository.findRoleAuthorities(r.role.ID.String())

	if err != nil {
		return nil, err
	}

	var resolvers []*authorityResolver

	for _, authority := range authorities {
		resolvers = append(resolvers, &authorityResolver{r.repository, authority})
	}

	return resolvers, nil
}

func (r *roleResolver) Users() ([]*userResolver, error) {
	return nil, nil
}
