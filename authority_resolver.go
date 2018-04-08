package main

import (
	"github.com/graph-gophers/graphql-go"
)

type authorityResolver struct {
	repository repository
	authority  *authority
}

func (r *authorityResolver) ID() graphql.ID {
	return graphQLID(r.authority.ID)
}

func (r *authorityResolver) Version() int32 {
	return int32(r.authority.Version)
}

func (r *authorityResolver) Name() string {
	return r.authority.Name
}

func (r *authorityResolver) Events() ([]*eventResolver, error) {
	return nil, nil
}

func (r *authorityResolver) Roles() ([]*roleResolver, error) {
	return nil, nil
}
