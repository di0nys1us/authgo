package main

import (
	graphql "github.com/neelance/graphql-go"
)

type authorityResolver struct {
	repository repository
	a          *authority
}

func (r *authorityResolver) ID() graphql.ID {
	return intToID(r.a.ID)
}

func (r *authorityResolver) Version() int32 {
	return int32(r.a.Version)
}

func (r *authorityResolver) Name() string {
	return r.a.Name
}

func (r *authorityResolver) Events() ([]*eventResolver, error) {
	return nil, nil
}

func (r *authorityResolver) Roles() ([]*roleResolver, error) {
	return nil, nil
}
