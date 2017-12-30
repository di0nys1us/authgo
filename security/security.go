package security

import (
	"context"
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
)

var (
	TimeFunc                = time.Now
	errInvalidSigningMethod = errors.New("authgo/security: invalid signing method")
	errMissingSecurityKey   = errors.New("authgo/security: missing environment variable AUTHGO_SECURITY_KEY")
)

type Subject interface {
	ID() int
	Username() string
	Password() string
	Administrator() bool
	Enabled() bool
}

type user struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Administrator bool   `json:"administrator"`
}

type jwtClaimsKey string

type authentication struct {
	subject     Subject
	tokenHolder *tokenHolder
}

type authorization struct {
	claims *jwtClaims
}

type findSubjectFunc func(email string) (Subject, error)

type defaultSecurity struct {
	findSubjectFunc findSubjectFunc
}

func NewSecurity(findSubjectFunc findSubjectFunc) *defaultSecurity {
	return &defaultSecurity{findSubjectFunc}
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

	user := &user{a.subject.ID(), a.subject.Username(), a.subject.Administrator()}

	return httpgo.WriteJSON(w, http.StatusOK, user)
}

func (s *defaultSecurity) authenticateRequest(r *http.Request) (*authentication, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when reading credentials")
	}

	subject, err := s.resolveSubject(r.Form.Get("email"), r.Form.Get("password"))

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when resolving subject")
	}

	t, err := createToken(subject)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when creating token")
	}

	return &authentication{subject, t}, nil
}

func (s *defaultSecurity) resolveSubject(email, password string) (Subject, error) {
	if s.findSubjectFunc == nil {
		return nil, errors.New("authgo/security: findSubjectFunc is not set")
	}

	subject, err := s.findSubjectFunc(email)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when finding subject")
	}

	if subject == nil {
		return nil, errors.New("authgo/security: subject is nil")
	}

	if !subject.Enabled() {
		return nil, errors.New("authgo/security: subject is not enabled")
	}

	err = validateHashedPassword(subject.Password(), password)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: invalid password")
	}

	return subject, nil
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
