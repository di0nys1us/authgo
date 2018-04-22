package security

import (
	"context"
	"html/template"
	"net/http"

	"github.com/di0nys1us/httpgo"
	"github.com/pkg/errors"
)

func (s *security) Authenticate(w http.ResponseWriter, r *http.Request) error {
	authN, err := s.authenticateRequest(r)

	if err != nil {
		return errors.WithStack(httpgo.ErrorWithStatusCode(http.StatusUnauthorized, err))
	}

	setAuthenticationCookie(w, authN)

	return httpgo.WriteJSON(w, http.StatusOK, nil)
}

func Authorize(next http.Handler) http.Handler {
	return httpgo.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		authZ, err := authorizeRequest(r)

		if err != nil {
			return errors.WithStack(err) // TODO 403
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxKeyUserID, authZ.jwtClaims.UserID)
		ctx = context.WithValue(ctx, ctxKeyUserEmail, authZ.jwtClaims.Subject)

		next.ServeHTTP(w, r.WithContext(ctx))

		return nil
	})
}

func GetLogin(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./templates/login.html")

	if err != nil {
		return errors.WithStack(err)
	}

	return tmpl.Execute(w, nil)
}

func (s *security) PostLogin(w http.ResponseWriter, r *http.Request) error {
	authN, err := s.authenticateRequest(r)

	if err != nil {
		return errors.WithStack(err)
	}

	setAuthenticationCookie(w, authN)

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, logoutCookie)

	return nil
}
