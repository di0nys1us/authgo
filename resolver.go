package authgo

import (
	"log"
	"strconv"

	"github.com/neelance/graphql-go"
)

func intToID(v int) graphql.ID {
	return graphql.ID(strconv.Itoa(v))
}

type resolver struct {
	db *db
}

func (r *resolver) Users() ([]*userResolver, error) {
	tx, err := r.db.begin()

	if err != nil {
		return nil, err
	}

	users, err := tx.findAllUsers()

	if err != nil {
		return nil, err
	}

	defer func() {
		err := tx.Commit()

		if err != nil {
			log.Print(err)
		}
	}()

	var resolvers []*userResolver

	for _, user := range users {
		resolvers = append(resolvers, &userResolver{user})
	}

	return resolvers, nil
}

type userResolver struct {
	u *user
}

func (r *userResolver) ID() graphql.ID {
	return intToID(r.u.ID)
}

func (r *userResolver) FirstName() string {
	return r.u.FirstName
}

func (r *userResolver) LastName() string {
	return r.u.LastName
}
