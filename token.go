package main

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type jwtClaims struct {
	*jwt.StandardClaims
	UserID uuid.UUID `json:"uid"`
}

type tokenHolder struct {
	token       *jwt.Token
	signedToken string
	claims      *jwtClaims
	expiresAt   time.Time
}

func createToken(user *user) (*tokenHolder, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when creating uuid based id")
	}

	now := timeFunc()
	expiresAt := now.AddDate(0, 0, 1)
	claims := &jwtClaims{
		&jwt.StandardClaims{
			Audience:  jwtAudience,
			ExpiresAt: expiresAt.Unix(),
			Id:        id.String(),
			IssuedAt:  now.Unix(),
			Issuer:    jwtIssuer,
			NotBefore: now.Unix(),
			Subject:   user.Email,
		},
		user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	securityKey, err := resolveSecurityKey(token)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when resolving security key")
	}

	signedToken, err := token.SignedString(securityKey)

	if err != nil {
		return nil, errors.Wrap(err, "authgo: error when signing the token")
	}

	return &tokenHolder{token, signedToken, claims, expiresAt}, nil
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
