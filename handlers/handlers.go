package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/di0nys1us/authgo/repository"
	"github.com/di0nys1us/httpgo"
)

// Handler TODO
type Handler struct {
	Repository *repository.Repository
}

// NewHandler TODO
func NewHandler(r *repository.Repository) *Handler {
	return &Handler{r}
}

// GetUsers TODO
func (a *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := a.Repository.FindUsers()

	if err != nil {
		return err
	}

	return httpgo.WriteJSON(w, http.StatusOK, users)
}

// GetUser TODO
func (a *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user, err := a.Repository.FindUser(id)

	if err != nil {
		return err
	}

	if user == nil {
		return httpgo.WriteJSON(w, http.StatusNotFound, nil)
	}

	return httpgo.WriteJSON(w, http.StatusOK, user)
}

// PostUser TODO
func (a *Handler) PostUser(w http.ResponseWriter, r *http.Request) error {
	user := &repository.User{}
	err := httpgo.ReadJSON(r.Body, user)

	if err != nil {
		return err
	}

	err = a.Repository.CreateUser(user)

	if err != nil {
		return err
	}

	return httpgo.WriteJSON(w, http.StatusCreated, user)
}

// PutUser TODO
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

	err = a.Repository.UpdateUser(user)

	if err != nil {
		return err
	}

	return httpgo.WriteJSON(w, http.StatusOK, user)
}
