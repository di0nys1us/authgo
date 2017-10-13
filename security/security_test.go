package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dgrijalva/jwt-go"
	"github.com/di0nys1us/httpgo"
)

func TestGetClaimsFromContext(t *testing.T) {
	testClaims := &jwtClaims{
		StandardClaims: &jwt.StandardClaims{Subject: "test@test.net"},
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

func TestAuthorize(t *testing.T) {
	handler := Authorize(httpgo.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	handler.ServeHTTP(w, r)

	assert.Exactly(t, http.StatusUnauthorized, w.Code)
}

func TestInvalidate(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	err := Invalidate(w, r)

	assert.Nil(t, err)
}
