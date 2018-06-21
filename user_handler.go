package main

import (
	"net/http"

	"github.com/di0nys1us/authgo/model"

	"github.com/di0nys1us/authgo/sqlgo"

	"github.com/go-chi/chi"

	"github.com/di0nys1us/httpgo"
	"github.com/pkg/errors"
)

type userHandler struct {
	db *sqlgo.DB
}

func (h *userHandler) getUsers(w http.ResponseWriter, r *http.Request) error {
	email := r.URL.Query().Get("email")

	if email != "" {
		user := &model.User{Email: email}

		tx, err := h.db.Start()

		if err != nil {
			return err
		}

		err = user.Find(tx)

		if err != nil {
			return errors.WithStack(err)
		}

		err = tx.Commit()

		if err != nil {
			return err
		}

		return httpgo.WriteJSON(w, http.StatusOK, user)
	}

	users := &model.Users{}

	err := users.FindAll(nil)

	if err != nil {
		return errors.WithStack(err)
	}

	return httpgo.WriteJSON(w, http.StatusOK, users)
}

func (h *userHandler) getUser(w http.ResponseWriter, r *http.Request) error {
	userID := chi.URLParam(r, "userID")

	user := &model.User{ID: userID}

	tx, err := h.db.Start()

	if err != nil {
		return err
	}

	err = user.Find(tx)

	if err != nil {
		return errors.WithStack(err)
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return httpgo.WriteJSON(w, http.StatusOK, user)
}
