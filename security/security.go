package security

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	environmentSecurityKey = "AUTHGO_SECURITY_KEY"
	jwtAudience            = "eies.land"
	jwtIssuer              = "authgo"
	jwtCookieName          = "authgo_token"
	ctxKeyUserID           = contextKeyUserID("ctxKeyUserID")
	ctxKeyUserEmail        = contextKeyUserEmail("ctxKeyUserEmail")
	formKeyEmail           = "email"
	formKeyPassword        = "password"
	UnknownUserID          = "UnknownUserID"
	UnknownUserEmail       = "UnknownUserEmail"
)

var (
	TimeFunc                = time.Now
	errInvalidSigningMethod = errors.New("authgo: invalid signing method")
	errMissingSecurityKey   = errors.New("authgo: missing environment variable AUTHGO_SECURITY_KEY")
	logoutCookie            = &http.Cookie{
		Name:     jwtCookieName,
		Value:    "",
		Expires:  time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		HttpOnly: true,
		Secure:   false,
	}
)

type contextKey string
type contextKeyUserID contextKey
type contextKeyUserEmail contextKey

type Subject interface {
	UserID() string
	UserEmail() string
	UserPassword() string
	UserActive() bool
}

type authentication struct {
	jwtToken *jwtToken
	subj     Subject
}

type authorization struct {
	jwtClaims *jwtClaims
}

type subjectByEmailFinder interface {
	FindSubjectByEmail(email string) (Subject, error)
}

type security struct {
	subjectByEmailFinder
}

func New(subjectFinder subjectByEmailFinder) *security {
	return &security{subjectFinder}
}

func UserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(ctxKeyUserID).(string); ok && userID != "" {
		return userID
	}

	return UnknownUserID
}

func UserEmailFromContext(ctx context.Context) string {
	if userEmail, ok := ctx.Value(ctxKeyUserEmail).(string); ok && userEmail != "" {
		return userEmail
	}

	return UnknownUserEmail
}

func (s *security) authenticateRequest(r *http.Request) (*authentication, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	subj, err := s.resolveSubject(r.Form.Get(formKeyEmail), r.Form.Get(formKeyPassword))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	jwtToken, err := createToken(subj)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &authentication{jwtToken, subj}, nil
}

func (s *security) resolveSubject(email, password string) (Subject, error) {
	subj, err := s.FindSubjectByEmail(email)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if subj == nil {
		return nil, errors.New("authgo: subject is nil")
	}

	if !subj.UserActive() {
		return nil, errors.New("authgo: user is not active")
	}

	err = validateHashedPassword(subj.UserPassword(), password)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return subj, nil
}

func authorizeRequest(r *http.Request) (*authorization, error) {
	cookie, err := r.Cookie(jwtCookieName)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, resolveSecurityKey)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !token.Valid {
		return nil, errors.New("authgo: invalid token")
	}

	jwt.TimeFunc = TimeFunc

	if err := claims.Valid(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &authorization{claims}, nil
}

func setAuthenticationCookie(w http.ResponseWriter, authN *authentication) {
	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    authN.jwtToken.signedToken,
		Expires:  authN.jwtToken.expiresAt,
		HttpOnly: true,
		Secure:   false,
	})
}
