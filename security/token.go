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
	UserID string `json:"uid"`
}

type jwtToken struct {
	signedToken string
	claims      *jwtClaims
	expiresAt   time.Time
}

func createToken(subj Subject) (*jwtToken, error) {
	id, err := uuid.NewV1()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	now := TimeFunc()
	expiresAt := now.AddDate(0, 0, 1)
	claims := &jwtClaims{
		&jwt.StandardClaims{
			Audience:  jwtAudience,
			ExpiresAt: expiresAt.Unix(),
			Id:        id.String(),
			IssuedAt:  now.Unix(),
			Issuer:    jwtIssuer,
			NotBefore: now.Unix(),
			Subject:   subj.UserEmail(),
		},
		subj.UserID(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	securityKey, err := resolveSecurityKey(token)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	signedToken, err := token.SignedString(securityKey)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &jwtToken{signedToken, claims, expiresAt}, nil
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
