package main

import (
	"strconv"

	"github.com/neelance/graphql-go"
)

func intToID(v int) graphql.ID {
	return graphql.ID(strconv.Itoa(v))
}

type rootResolver struct {
	*rootMutation
	repository repository
}

func (r *rootResolver) Users() ([]*userResolver, error) {
	users, err := r.repository.findAllUsers()

	if err != nil {
		return nil, err
	}

	var resolvers []*userResolver

	for _, user := range users {
		resolvers = append(resolvers, &userResolver{r.repository, user})
	}

	return resolvers, nil
}

func (r *rootResolver) User(args struct {
	ID    *graphql.ID
	Email *string
}) (*userResolver, error) {
	var user *user
	var err error

	if args.ID != nil {
		user, err = r.repository.findUserByID(string(*args.ID))
	}

	if args.Email != nil {
		user, err = r.repository.findUserByEmail(*args.Email)
	}

	if err != nil {
		return nil, err
	}

	return &userResolver{r.repository, user}, nil
}

type userResolver struct {
	repository repository
	u          *user
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
	events, err := r.repository.findEventsByUserID(strconv.Itoa(r.u.ID))

	if err != nil {
		return nil, err
	}

	var resolvers []*eventResolver

	for _, event := range events {
		resolvers = append(resolvers, &eventResolver{r.repository, event})
	}

	return resolvers, nil
}

func (r *userResolver) Roles() ([]*roleResolver, error) {
	roles, err := r.repository.findUserRoles(strconv.Itoa(r.u.ID))

	if err != nil {
		return nil, err
	}

	var resolvers []*roleResolver

	for _, role := range roles {
		resolvers = append(resolvers, &roleResolver{r.repository, role})
	}

	return resolvers, nil
}

type eventResolver struct {
	repository repository
	event      *event
}

func (r *eventResolver) ID() graphql.ID {
	return intToID(r.event.ID)
}

func (r *eventResolver) CreatedBy() (*userResolver, error) {
	user, err := r.repository.findUserByID(strconv.Itoa(r.event.CreatedBy))

	if err != nil {
		return nil, err
	}

	return &userResolver{r.repository, user}, nil
}

func (r *eventResolver) CreatedAt() string {
	return r.event.CreatedAt.String()
}

func (r *eventResolver) Type() (string, error) {
	return r.event.Type, nil
}

func (r *eventResolver) Description() string {
	return r.event.Description
}

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
