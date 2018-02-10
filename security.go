package authgo

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
	keyEmail               = "email"
	keyPassword            = "password"
)

var (
	timeFunc                = time.Now
	errInvalidSigningMethod = errors.New("authgo: invalid signing method")
	errMissingSecurityKey   = errors.New("authgo: missing environment variable AUTHGO_SECURITY_KEY")
)

type jwtClaimsKey string

type authentication struct {
	user        *user
	tokenHolder *tokenHolder
}

type authorization struct {
	claims *jwtClaims
}

type security struct {
	db *db
}

func newSecurity(db *db) *security {
	return &security{db}
}

func newContextWithClaims(c context.Context, claims *jwtClaims) context.Context {
	return context.WithValue(c, claimsKey, claims)
}

func getClaimsFromContext(c context.Context) (*jwtClaims, bool) {
	claims, ok := c.Value(claimsKey).(*jwtClaims)
	return claims, ok
}

func (sec *security) authenticate(w http.ResponseWriter, r *http.Request) error {
	aut, err := sec.authenticateRequest(r)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    aut.tokenHolder.signedToken,
		Expires:  aut.tokenHolder.expires,
		HttpOnly: true,
		Secure:   false,
	})

	return httpgo.WriteJSON(w, http.StatusOK, aut.user)
}

func (sec *security) authenticateRequest(r *http.Request) (*authentication, error) {
	err := r.ParseForm()

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when reading credentials")
	}

	subject, err := sec.resolveUser(r.Form.Get(keyEmail), r.Form.Get(keyPassword))

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when resolving subject")
	}

	t, err := createToken(subject)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when creating token")
	}

	return &authentication{subject, t}, nil
}

func (sec *security) resolveUser(email, password string) (*user, error) {
	user, err := sec.db.findUserByEmail(email)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when finding user")
	}

	if user == nil {
		return nil, errors.New("authgo: user is nil")
	}

	if !user.Enabled {
		return nil, errors.New("authgo: user is not enabled")
	}

	err = validateHashedPassword(user.Password, password)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: invalid password")
	}

	return user, nil
}

func authorize(next http.Handler) http.Handler {
	return httpgo.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		auth, err := authorizeRequest(r)

		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(newContextWithClaims(r.Context(), auth.claims)))

		return nil
	})
}

func authorizeRequest(r *http.Request) (*authorization, error) {
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

	return &authorization{claims}, nil
}

func (sec *security) getLogin(w http.ResponseWriter, r *http.Request) error {
	return writeString(w, templateLogin)
}

func (sec *security) postLogin(w http.ResponseWriter, r *http.Request) error {
	aut, err := sec.authenticateRequest(r)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     jwtCookieName,
		Value:    aut.tokenHolder.signedToken,
		Expires:  aut.tokenHolder.expires,
		HttpOnly: true,
		Secure:   false,
	})

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
