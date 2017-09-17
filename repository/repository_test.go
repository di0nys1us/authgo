// +build integration

package repository

import (
	"flag"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

var (
	repository Repository
	now        = time.Now()
	testUser1  = &User{
		ID:            1,
		Version:       0,
		Deleted:       false,
		CreatedAt:     now,
		CreatedBy:     "it",
		ModifiedAt:    now,
		ModifiedBy:    "it",
		FirstName:     "Bar",
		LastName:      "Foo",
		Email:         "bar@foo.land",
		Password:      "password",
		Enabled:       true,
		Administrator: true,
	}
	testUser2 = &User{
		ID:            2,
		Version:       0,
		Deleted:       false,
		CreatedAt:     now,
		CreatedBy:     "it",
		ModifiedAt:    now,
		ModifiedBy:    "it",
		FirstName:     "Foo",
		LastName:      "Bar",
		Email:         "foo@bar.land",
		Password:      "password",
		Enabled:       true,
		Administrator: false,
	}
)

func createRepository() Repository {
	r, err := NewRepository()

	if err != nil {
		log.Fatalln(err)
	}

	return r
}

func TestMain(m *testing.M) {
	flag.Parse()

	repository = createRepository()
	defer repository.Close()

	os.Exit(m.Run())
}

func TestFindUsers(t *testing.T) {
	users, err := repository.FindUsers()

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	if n := len(users); n != 2 {
		t.Fatalf("got %v, want %v", n, 2)
	}
}

func TestFindUser(t *testing.T) {
	user, err := repository.FindUser(strconv.Itoa(testUser1.ID))

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	if !reflect.DeepEqual(user, testUser1) {
		t.Errorf("got %v, want %v", user, testUser1)
	}
}

func TestFindUserByEmail(t *testing.T) {
	user, err := repository.FindUserByEmail(testUser2.Email)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}

	if !reflect.DeepEqual(user, testUser2) {
		t.Errorf("got %v, want %v", user, testUser2)
	}
}

func TestSaveUser(t *testing.T) {
	err := repository.SaveUser(testUser1)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}

func TestUpdateUser(t *testing.T) {
	err := repository.UpdateUser(testUser2)

	if err != nil {
		t.Fatalf("got %v, want nil", err)
	}
}
