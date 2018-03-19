package main

import (
	"strconv"

	graphql "github.com/neelance/graphql-go"
)

type roleResolver struct {
	repository repository
	r          *role
}

func (r *roleResolver) ID() graphql.ID {
	return intToID(r.r.ID)
}

func (r *roleResolver) Version() int32 {
	return int32(r.r.Version)
}

func (r *roleResolver) Name() string {
	return r.r.Name
}

func (r *roleResolver) Events() ([]*eventResolver, error) {
	return nil, nil
}

func (r *roleResolver) Authorities() ([]*authorityResolver, error) {
	authorities, err := r.repository.findRoleAuthorities(strconv.Itoa(r.r.ID))

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
