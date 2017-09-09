package handlers

import (
	"testing"

	"github.com/di0nys1us/authgo/repository"
)

type MockRepository struct {
	User *repository.User
}

func TestGetUsers(t *testing.T) {
}
