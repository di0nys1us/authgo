package security

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func validateHashedPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return errors.WithStack(err)
}

func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), errors.WithStack(err)
}
