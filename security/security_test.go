package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/di0nys1us/httpgo"
)

func TestResolveSecurityKey(t *testing.T) {
	os.Setenv("", "")

	token := jwt.New(jwt.SigningMethodES256)

	securityKey, err := resolveSecurityKey(token)

	if securityKey != nil {
		t.Errorf("got %v, want nil", securityKey)
	}

	if _, ok := err.(ErrUnexpectedSigningMethod); !ok {
		t.Errorf("got %T, want %T", err, ErrUnexpectedSigningMethod(""))
	}

	token = jwt.New(jwt.SigningMethodHS256)

	securityKey, err = resolveSecurityKey(token)

	if securityKey != nil {
		t.Errorf("got %v, want nil", securityKey)
	}

	if err != ErrMissingSecurityKey {
		t.Errorf("got %s, want %s", err, ErrMissingSecurityKey)
	}
}

func TestGetClaimsFromContext(t *testing.T) {
	testClaims := &JWTClaims{
		StandardClaims: jwt.StandardClaims{Subject: "foo@bar.net"},
	}

	c := context.Background()

	if claims, ok := GetClaimsFromContext(c); ok {
		t.Errorf("got %v, want nil", claims)
	}

	c = newContextWithClaims(c, testClaims)

	if claims, ok := GetClaimsFromContext(c); !ok {
		t.Errorf("got %v, want %v", claims, testClaims)
	}
}

func TestValidateJWT(t *testing.T) {
	handler := ValidateJWT(httpgo.ResponseHandlerFunc(func(w http.ResponseWriter, r *http.Request) (*httpgo.Response, error) {
		return httpgo.ResponseOK().WithBody(true), nil
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add(httpgo.HeaderAuthorization, "Bearer TEST")

	handler.ServeHTTP(w, r)

	if code := w.Result().StatusCode; code != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", code, http.StatusUnauthorized)
	}
}
