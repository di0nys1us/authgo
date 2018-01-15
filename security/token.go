package security

import (
	"os"
	"time"

	"github.com/di0nys1us/authgo/domain"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type jwtClaims struct {
	*jwt.StandardClaims
	UserID int `json:"uid"`
}

type tokenHolder struct {
	token       *jwt.Token
	signedToken string
	claims      *jwtClaims
	expires     time.Time
}

func createToken(user *domain.User) (*tokenHolder, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when creating uuid based id")
	}

	now := TimeFunc()
	expires := now.AddDate(0, 0, 1)
	claims := &jwtClaims{
		&jwt.StandardClaims{
			Audience:  jwtAudience,
			ExpiresAt: expires.Unix(),
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
		return nil, errors.Wrap(err, "authgo/security: error when resolving security key")
	}

	signedToken, err := token.SignedString(securityKey)

	if err != nil {
		return nil, errors.Wrap(err, "authgo/security: error when signing the token")
	}

	return &tokenHolder{token, signedToken, claims, expires}, nil
}

func resolveSecurityKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errInvalidSigningMethod
	}

	securityKey, ok := os.LookupEnv(environmentSecurityKey)

	if !ok {
		return nil, errMissingSecurityKey
	}

	return []byte(securityKey), nil
}
