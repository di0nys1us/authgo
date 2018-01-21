package security

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/di0nys1us/authgo/domain"

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
	TimeFunc                = time.Now
	errInvalidSigningMethod = errors.New("authgo/security: invalid signing method")
	errMissingSecurityKey   = errors.New("authgo/security: missing environment variable AUTHGO_SECURITY_KEY")
)

type jwtClaimsKey string

type authentication struct {
	user        *domain.User
	tokenHolder *tokenHolder
}

type authorization struct {
	claims *jwtClaims
}

type findUserFunc func(email string) (*domain.User, error)

type defaultSecurity struct {
	findUserFunc findUserFunc
}

func NewSecurity(findUserFunc findUserFunc) *defaultSecurity {
	return &defaultSecurity{findUserFunc}
}

func newContextWithClaims(c context.Context, claims *jwtClaims) context.Context {
	return context.WithValue(c, claimsKey, claims)
}

func GetClaimsFromContext(c context.Context) (*jwtClaims, bool) {
	claims, ok := c.Value(claimsKey).(*jwtClaims)
	return claims, ok
}

func (s *defaultSecurity) Authenticate(w http.ResponseWriter, r *http.Request) error {
	a, err := s.authenticateRequest(r)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    a.tokenHolder.signedToken,
		Expires:  a.tokenHolder.expires,
		HttpOnly: true,
		Secure:   false,
	})

	return httpgo.WriteJSON(w, http.StatusOK, a.user)
}

func (s *defaultSecurity) authenticateRequest(r *http.Request) (*authentication, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when reading credentials")
	}

	subject, err := s.resolveUser(r.Form.Get(keyEmail), r.Form.Get(keyPassword))

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when resolving subject")
	}

	t, err := createToken(subject)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when creating token")
	}

	return &authentication{subject, t}, nil
}

func (s *defaultSecurity) resolveUser(email, password string) (*domain.User, error) {
	if s.findUserFunc == nil {
		return nil, errors.New("authgo/security: findUserFunc is not set")
	}

	user, err := s.findUserFunc(email)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when finding user")
	}

	if user == nil {
		return nil, errors.New("authgo/security: user is nil")
	}

	if !user.Enabled {
		return nil, errors.New("authgo/security: user is not enabled")
	}

	err = validateHashedPassword(user.Password, password)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: invalid password")
	}

	return user, nil
}

func Authorize(next http.Handler) http.Handler {
	return httpgo.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		a, err := authorizeRequest(r)

		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(newContextWithClaims(r.Context(), a.claims)))

		return nil
	})
}

func authorizeRequest(r *http.Request) (*authorization, error) {
	cookie, err := r.Cookie(jwtCookieName)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: missing cookie")
	}

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, resolveSecurityKey)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: token error")
	}

	if !token.Valid {
		return nil, errors.Wrap(err, "authgo/security: invalid token")
	}

	jwt.TimeFunc = TimeFunc

	if err := claims.Valid(); err != nil {
		return nil, errors.Wrap(err, "authgo/security: invalid claims")
	}

	return &authorization{claims}, nil
}

func (s *defaultSecurity) GetLogin(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("template/authenticate.html")

	if err != nil {
		return err
	}

	return tmpl.Execute(w, nil)
}

func (s *defaultSecurity) PostLogin(w http.ResponseWriter, r *http.Request) error {
	a, err := s.authenticateRequest(r)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    a.tokenHolder.signedToken,
		Expires:  a.tokenHolder.expires,
		HttpOnly: true,
		Secure:   false,
	})

	http.Redirect(w, r, "/graphql", http.StatusSeeOther)

	return nil
}

func validateHashedPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return errors.Wrap(err, "authgo/security: error when validating hashed password")
}

func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), errors.Wrap(err, "authgo/security: error when generating hashed password")
}

func Invalidate(w http.ResponseWriter, r *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    "",
		Expires:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		HttpOnly: true,
		Secure:   false,
	})

	return nil
}
