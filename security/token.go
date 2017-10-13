package security

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type jwtClaims struct {
	*jwt.StandardClaims
	UserID        int  `json:"uid"`
	Administrator bool `json:"adm"`
}

type tokenHolder struct {
	token       *jwt.Token
	signedToken string
	claims      *jwtClaims
	expires     time.Time
}

func createToken(s Subject) (*tokenHolder, error) {
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

	return &tokenHolder{t, signedToken, claims, expires}, nil
}

func resolveSecurityKey(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errInvalidSigningMethod
	}

	securityKey, ok := os.LookupEnv(environmentSecurityKey)

	if !ok {
		return nil, errMissingSecurityKey
	}

	return []byte(securityKey), nil
}
