package main

import (
	"strconv"

	"github.com/neelance/graphql-go"
)

func intToID(v int) graphql.ID {
	return graphql.ID(strconv.Itoa(v))
}

type rootResolver struct {
	*rootQuery
	*rootMutation
}
