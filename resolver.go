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
	var users []*user
	var err error

	err = r.db.do(func(tx *tx) {
		users, err = tx.findAllUsers()
	})

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

func (r *userResolver) Events() ([]*eventResolver, error) {
	return nil, nil
}

func (r *userResolver) Roles() ([]*roleResolver, error) {
	return nil, nil
}

type eventResolver struct {
	evt *event
}

func (r *eventResolver) ID() graphql.ID {
	return intToID(r.evt.ID)
}

func (r *eventResolver) CreatedBy() (*userResolver, error) {
	return nil, nil
}

func (r *eventResolver) CreatedAt() string {
	return r.evt.CreatedAt.String()
}

func (r *eventResolver) Type() (string, error) {
	return "TODO", nil
}

func (r *eventResolver) Description() string {
	return r.evt.Description
}

type roleResolver struct {
	r *role
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
	return nil, nil
}

func (r *roleResolver) Users() ([]*userResolver, error) {
	return nil, nil
}

type authorityResolver struct {
	a *authority
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
