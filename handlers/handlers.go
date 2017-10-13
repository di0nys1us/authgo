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

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{r}
}

func (a *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := a.Repository.FindUsers()

	if err != nil {
		return err
	}

	return httpgo.WriteJSON(w, http.StatusOK, users)
}

func (a *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user, err := a.Repository.FindUser(id)

	if err != nil {
		return err
	}

	if user == nil {
		return httpgo.WriteJSON(w, http.StatusNoContent, nil)
	}

	return httpgo.WriteJSON(w, http.StatusOK, user)
}

func (a *Handler) PostUser(w http.ResponseWriter, r *http.Request) error {
	user := &repository.User{}
	err := httpgo.ReadJSON(r.Body, user)

	if err != nil {
		return err
	}

	err = prepareUser(user, r.Context())

	if err != nil {
		return err
	}

	err = a.Repository.SaveUser(user)

	if err != nil {
		return err
	}

	user.Password = ""

	return httpgo.WriteJSON(w, http.StatusCreated, user)
}

func (a *Handler) PutUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	user := &repository.User{}
	err := httpgo.ReadJSON(r.Body, user)

	if err != nil {
		return err
	}

	if id != strconv.Itoa(user.ID) {
		return errors.New("authgo/handlers: id mismatch")
	}

	err = prepareUser(user, r.Context())

	if err != nil {
		return err
	}

	err = a.Repository.UpdateUser(user)

	if err != nil {
		return err
	}

	user.Password = ""

	return httpgo.WriteJSON(w, http.StatusOK, user)
}

func prepareUser(u *repository.User, ctx context.Context) error {
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
