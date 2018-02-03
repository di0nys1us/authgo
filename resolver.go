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

	defer func() {
		err := tx.Commit()

		if err != nil {
			log.Print(err)
		}
	}()

	users, err := tx.findAllUsers()

	if err != nil {
		return nil, err
	}

	var resolvers []*userResolver

	for _, user := range users {
		resolvers = append(resolvers, &userResolver{user})
	}

	return resolvers, nil
}

func (r *resolver) User(args struct {
	ID    *graphql.ID
	Email *string
}) (*userResolver, error) {
	tx, err := r.db.begin()

	if err != nil {
		return nil, err
	}

	defer func() {
		err := tx.Commit()

		if err != nil {
			log.Print(err)
		}
	}()

	var user *user

	if args.ID != nil {
		user, err = tx.findUserByID(string(*args.ID))
	}

	if args.Email != nil {
		user, err = tx.findUserByEmail(*args.Email)
	}

	if err != nil {
		return nil, err
	}

	return &userResolver{user}, nil
}

type userResolver struct {
	u *user
}

func (r *userResolver) ID() graphql.ID {
	return intToID(r.u.ID)
}

func (r *userResolver) Version() int32 {
	return int32(r.u.Version)
}

func (r *userResolver) FirstName() string {
	return r.u.FirstName
}

func (r *userResolver) LastName() string {
	return r.u.LastName
}

func (r *userResolver) Email() string {
	return r.u.Email
}

func (r *userResolver) Password() string {
	return r.u.Password
}

func (r *userResolver) Enabled() bool {
	return r.u.Enabled
}

func (r *userResolver) Deleted() bool {
	return r.u.Deleted
}
