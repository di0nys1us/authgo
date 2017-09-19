package security

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/di0nys1us/httpgo"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtAudience   = "eies.land"
	jwtIssuer     = "authgo"
	jwtCookieName = "authgo_token"
	claimsKey     = jwtClaimsKey("jwtClaimsKey")
	bearerToken   = "Bearer %s"
)

var (
	TimeFunc = time.Now
)

type Subject interface {
	GetID() int
	GetEmail() string
	GetPassword() string
	IsAdministrator() bool
	IsInactive() bool
}

type user struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Administrator bool   `json:"administrator"`
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type jwtClaims struct {
	*jwt.StandardClaims
	UserID        int  `json:"uid"`
	Administrator bool `json:"adm"`
}

type jwtClaimsKey string

type token struct {
	token       *jwt.Token
	signedToken string
	claims      *jwtClaims
	expires     time.Time
}

type authentication struct {
	subject Subject
	token   *token
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

func (s *defaultSecurity) Authenticate(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
	a, err := s.authenticateRequest(r)

	if err != nil {
		log.Println(err)
		return httpgo.ResponseForbidden(), nil
	}

	cookie := http.Cookie{
		Name:     jwtCookieName,
		Value:    a.token.signedToken,
		Expires:  a.token.expires,
		HttpOnly: true,
	}
	user := &user{a.subject.GetID(), a.subject.GetEmail(), a.subject.IsAdministrator()}

	return httpgo.ResponseOK().WithCookie(cookie).WithBody(user), nil
}

func (s *defaultSecurity) authenticateRequest(r *http.Request) (*authentication, error) {
	c := credentials{}
	err := httpgo.ReadJSON(r.Body, &c)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when reading credentials")
	}

	subject, err := s.resolveSubject(c)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when resolving subject")
	}

	t, err := createToken(subject)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when creating token")
	}

	return &authentication{subject, t}, nil
}

func (s *defaultSecurity) resolveSubject(c credentials) (Subject, error) {
	if s.findSubjectFunc == nil {
		return nil, errors.New("authgo/security: findSubjectFunc is not set")
	}

	subject, err := s.findSubjectFunc(c.Email)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when finding subject")
	}

	if subject == nil {
		return nil, errors.New("authgo/security: subject is nil")
	}

	if subject.IsInactive() {
		return nil, errors.New("authgo/security: subject is inactive")
	}

	err = validateHashedPassword(subject.GetPassword(), c.Password)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: invalid password")
	}

	return subject, nil
}

func createToken(s Subject) (*token, error) {
	now := TimeFunc()
	expires := now.AddDate(0, 0, 1)
	claims := &jwtClaims{
		&jwt.StandardClaims{
			Audience:  jwtAudience,
			ExpiresAt: expires.Unix(),
			Id:        uuid.NewV4().String(),
			IssuedAt:  now.Unix(),
			Issuer:    jwtIssuer,
			NotBefore: now.Unix(),
			Subject:   s.GetEmail(),
		},
		s.GetID(),
		s.IsAdministrator(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, err := resolveSecurityKey(t)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when resolving security key")
	}

	signedToken, err := t.SignedString(secretKey)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when signing the token")
	}

	return &token{t, signedToken, claims, expires}, nil
}

func Authorize(next http.Handler) http.Handler {
	return httpgo.ResponseHandlerFunc(func(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
		a, err := authorizeRequest(r)

		if err != nil {
			log.Println(err)
			return httpgo.ResponseUnauthorized(), nil
		}

		next.ServeHTTP(w, r.WithContext(newContextWithClaims(r.Context(), a.claims)))

		return nil, nil
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

func resolveSecurityKey(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("authgo/security: invalid signing method")
	}

	securityKey, ok := os.LookupEnv("AUTHGO_KEY")

	if !ok {
		return nil, errors.New("authgo/security: missing environment variable AUTHGO_KEY")
	}

	return []byte(securityKey), nil
}

func validateHashedPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return errors.Wrap(err, "authgo/security: error when validating hashed password")
}

func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), errors.Wrap(err, "authgo/security: error when generating hashed password")
}
