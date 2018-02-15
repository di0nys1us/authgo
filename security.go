package main

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/di0nys1us/httpgo"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtAudience            = "eies.land"
	jwtIssuer              = "authgo"
	jwtCookieName          = "authgo_token"
	claimsKey              = jwtClaimsKey("jwtClaimsKey")
	environmentSecurityKey = "AUTHGO_SECURITY_KEY"
	keyEmail               = "email"
	keyPassword            = "password"
)

var (
	timeFunc                = time.Now
	errInvalidSigningMethod = errors.New("authgo: invalid signing method")
	errMissingSecurityKey   = errors.New("authgo: missing environment variable AUTHGO_SECURITY_KEY")
)

type jwtClaimsKey string

type authenticationHolder struct {
	user        *user
	tokenHolder *tokenHolder
}

type authorizationHolder struct {
	claims *jwtClaims
}

type security struct {
	userByEmailFinder userByEmailFinder
}

func newSecurity(userByEmailFinder userByEmailFinder) *security {
	return &security{userByEmailFinder}
}

func newContextWithClaims(c context.Context, claims *jwtClaims) context.Context {
	return context.WithValue(c, claimsKey, claims)
}

func getClaimsFromContext(c context.Context) (*jwtClaims, bool) {
	claims, ok := c.Value(claimsKey).(*jwtClaims)
	return claims, ok
}

func (s *security) authenticate(w http.ResponseWriter, r *http.Request) error {
	authentication, err := s.authenticateRequest(r)

	if err != nil {
		return err
	}

	setAuthenticationCookie(w, authentication)

	return httpgo.WriteJSON(w, http.StatusOK, authentication.user)
}

func (s *security) authenticateRequest(r *http.Request) (*authenticationHolder, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when reading credentials")
	}

	user, err := s.resolveUser(r.Form.Get(keyEmail), r.Form.Get(keyPassword))

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when resolving subject")
	}

	token, err := createToken(user)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when creating token")
	}

	return &authenticationHolder{user, token}, nil
}

func (s *security) resolveUser(email, password string) (*user, error) {
	user, err := s.userByEmailFinder.findUserByEmail(email)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding user")
	}

	if user == nil {
		return nil, errors.New("authgo: user is nil")
	}

	if !user.Enabled {
		return nil, errors.New("authgo: user is not enabled")
	}

	if user.Deleted {
		return nil, errors.New("authgo: user is deleted")
	}

	err = validateHashedPassword(user.Password, password)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: invalid password")
	}

	return user, nil
}

func authorize(next http.Handler) http.Handler {
	return httpgo.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		authorization, err := authorizeRequest(r)

		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(newContextWithClaims(r.Context(), authorization.claims)))

		return nil
	})
}

func authorizeRequest(r *http.Request) (*authorizationHolder, error) {
	cookie, err := r.Cookie(jwtCookieName)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: missing cookie")
	}

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, resolveSecurityKey)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: token error")
	}

	if !token.Valid {
		return nil, errors.Wrap(err, "authgo: invalid token")
	}

	jwt.TimeFunc = timeFunc

	if err := claims.Valid(); err != nil {
		return nil, errors.Wrap(err, "authgo: invalid claims")
	}

	return &authorizationHolder{claims}, nil
}

func (s *security) getLogin(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./templates/login.html")

	if err != nil {
		return err
	}

	return tmpl.Execute(w, nil)
}

func (s *security) postLogin(w http.ResponseWriter, r *http.Request) error {
	authentication, err := s.authenticateRequest(r)

	if err != nil {
		return err
	}

	setAuthenticationCookie(w, authentication)

	http.Redirect(w, r, "/graphiql", http.StatusSeeOther)

	return nil
}

func validateHashedPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return errors.Wrap(err, "authgo: error when validating hashed password")
}

func generateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), errors.Wrap(err, "authgo: error when generating hashed password")
}

func setAuthenticationCookie(w http.ResponseWriter, a *authenticationHolder) {
	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    a.tokenHolder.signedToken,
		Expires:  a.tokenHolder.expiresAt,
		HttpOnly: true,
		Secure:   false,
	})
}

func invalidate(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    "",
		Expires:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		HttpOnly: true,
		Secure:   false,
	})

	return nil
}
