package authgo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	jwt "github.com/dgrijalva/jwt-go"
)

const testKey = "secret"

func TestResolveSecurityKey(t *testing.T) {
	token := jwt.New(jwt.SigningMethodES256)

	securityKey, err := resolveSecurityKey(token)

	assert.Nil(t, securityKey)
	assert.Exactly(t, errInvalidSigningMethod, err)

	token = jwt.New(jwt.SigningMethodHS256)

	securityKey, err = resolveSecurityKey(token)

	assert.Nil(t, securityKey)
	assert.Exactly(t, errMissingSecurityKey, err)

	os.Setenv(environmentSecurityKey, testKey)

	securityKey, err = resolveSecurityKey(token)

	assert.Exactly(t, []byte(testKey), securityKey)
	assert.Nil(t, err)
}
