package security

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
	jwtCookieName = "jwt"
	claimsKey     = JWTClaimsKey("JWTClaimsKey")
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

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FindSubjectFunc func(email string) (Subject, error)

type DefaultSecurity struct {
	FindSubjectFunc FindSubjectFunc
}

func NewSecurity(findSubjectFunc FindSubjectFunc) *DefaultSecurity {
	return &DefaultSecurity{findSubjectFunc}
}

type JWTClaims struct {
	jwt.StandardClaims
	UserID        int  `json:"uid"`
	Administrator bool `json:"adm"`
}

type JWTClaimsKey string

func newContextWithClaims(c context.Context, claims *JWTClaims) context.Context {
	return context.WithValue(c, claimsKey, claims)
}

func GetClaimsFromContext(c context.Context) (*JWTClaims, bool) {
	claims, ok := c.Value(claimsKey).(*JWTClaims)
	return claims, ok
}

func (s *DefaultSecurity) Authenticate(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
	c := Credentials{}
	err := httpgo.ReadJSON(r.Body, &c)

	if err != nil {
		log.Println(err)
		return httpgo.ResponseForbidden(), nil
	}

	subject, err := s.resolveSubject(c)

	if err != nil {
		log.Println(err)
		return httpgo.ResponseForbidden(), nil
	}

	token := createToken(subject)
	secretKey, err := resolveSecurityKey(token)

	if err != nil {
		log.Println(err)
		return httpgo.ResponseForbidden(), nil
	}

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		log.Println(err)
		return httpgo.ResponseForbidden(), nil
	}

	return httpgo.ResponseOK().WithHeader(httpgo.HeaderAuthorization, fmt.Sprintf(bearerToken, tokenString)), nil
}

func ValidateJWT(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
		bearerToken := r.Header.Get(httpgo.HeaderAuthorization)
		claims, err := extractClaims(bearerToken)

		if err != nil {
			log.Println(err)
			return httpgo.ResponseUnauthorized(), nil
		}

		next.ServeHTTP(w, r.WithContext(newContextWithClaims(r.Context(), claims)))

		return nil, nil
	}

	return httpgo.ResponseHandlerFunc(handler)
}

func extractClaims(bearerToken string) (*JWTClaims, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return nil, errors.New("authgo/security: missing bearer token")
	}

	tokenString := bearerToken[7:]
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, resolveSecurityKey)

	if err != nil || !token.Valid {
		return nil, errors.Wrap(err, "authgo/security: invalid token")
	}

	jwt.TimeFunc = TimeFunc

	if err := claims.Valid(); err != nil {
		return nil, errors.Wrap(err, "authgo/security: invalid claims")
	}

	return claims, nil
}

func (s *DefaultSecurity) resolveSubject(c Credentials) (Subject, error) {
	subject, err := s.FindSubjectFunc(c.Email)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when finding subject")
	}

	if subject == nil || subject.IsInactive() {
		return nil, errors.New("authgo/security: invalid subject")
	}

	err = validateHashedPassword(subject.GetPassword(), c.Password)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: invalid password")
	}

	return subject, nil
}

func createToken(s Subject) *jwt.Token {
	now := TimeFunc()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		jwt.StandardClaims{
			Audience:  jwtAudience,
			ExpiresAt: now.AddDate(0, 0, 1).Unix(),
			Id:        uuid.NewV4().String(),
			IssuedAt:  now.Unix(),
			Issuer:    jwtIssuer,
			NotBefore: now.Unix(),
			Subject:   s.GetEmail(),
		},
		s.GetID(),
		s.IsAdministrator(),
	})
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
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), errors.Wrap(err, "authgo/security: error when generating hashed password")
}
