package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/di0nys1us/authgo/repository"
	"github.com/di0nys1us/authgo/security"
	"github.com/di0nys1us/httpgo"
)

type UsersGetter interface {
	GetUsers(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error)
}

type UserGetter interface {
	GetUser(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error)
}

type UserPoster interface {
	PostUser(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error)
}

type UserPutter interface {
	PutUser(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error)
}

type API interface {
	UsersGetter
	UserGetter
	UserPoster
	UserPutter
}

type DefaultAPI struct {
	Repository repository.Repository
}

func NewAPI(repository repository.Repository) API {
	return &DefaultAPI{repository}
}

func (a *DefaultAPI) GetUsers(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
	users, err := a.Repository.FindUsers()

	if err != nil {
		return nil, err
	}

	return httpgo.ResponseOK().WithBody(users), nil
}

func (a *DefaultAPI) GetUser(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
	id := chi.URLParam(r, "id")

	user, err := a.Repository.FindUser(id)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return httpgo.ResponseNoContent(), nil
	}

	return httpgo.ResponseOK().WithBody(user), nil
}

func (a *DefaultAPI) PostUser(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
	user := &repository.User{}
	err := httpgo.ReadJSON(r.Body, user)

	if err != nil {
		return nil, err
	}

	err = updateUser(user, r.Context())

	if err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return httpgo.ResponseBadRequest().WithBody(err), nil
	}

	err = a.Repository.SaveUser(user)

	if err != nil {
		return nil, err
	}

	user.Password = ""

	return httpgo.ResponseOK().WithBody(user), nil
}

func (a *DefaultAPI) PutUser(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
	id := chi.URLParam(r, "id")
	user := &repository.User{}
	err := httpgo.ReadJSON(r.Body, user)

	if err != nil {
		return nil, err
	}

	if id != strconv.Itoa(user.ID) {
		return nil, errors.New("authgo/handlers: id mismatch")
	}

	err = updateUser(user, r.Context())

	if err != nil {
		return nil, err
	}

	err = a.Repository.UpdateUser(user)

	if err != nil {
		return nil, err
	}

	user.Password = ""

	return httpgo.ResponseOK().WithBody(user), nil
}

func updateUser(u *repository.User, ctx context.Context) error {
	var subject string

	if claims, ok := security.GetClaimsFromContext(ctx); ok {
		subject = claims.Subject
	} else {
		return errors.New("authgo/handlers: missing subject")
	}

	password, err := security.GenerateHashedPassword(u.Password)

	if err != nil {
		return err
	}

	u.CreatedBy = subject
	u.ModifiedBy = subject
	u.Password = password

	return nil
}
