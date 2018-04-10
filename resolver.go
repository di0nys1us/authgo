package main

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/satori/go.uuid"
)

func graphQLID(i interface{}) graphql.ID {
	switch v := i.(type) {
	case graphql.ID:
		return v
	case string:
		return graphql.ID(v)
	case uuid.UUID:
		return graphql.ID(v.String())
	default:
		return ""
	}
}

type rootResolver struct {
	*rootQuery
	*rootMutation
}
