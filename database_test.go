package main

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	db, err := newDB()

	defer db.DB.Close()

	if err != nil {
		t.Fatal(err)
	}

	uuid, err := db.generateUUID()

	if err != nil {
		t.Fatal(err)
	}

	t.Fatal(uuid)
}
