package database

import (
	"os"
	"testing"
)

func TestGetMembersWithCredit(t *testing.T) {
	os.Setenv("DB_CONNECTION_STRING", "postgresql://test:test@localhost:5432/membership")
	db, err := Setup()
	if err != nil {
		t.Error(err)
	}
	defer db.Release()

	members := db.GetMembersWithCredit()

	expectedMemberCount := 2
	actualMemberCount := len(members)
	if actualMemberCount != expectedMemberCount {
		t.Errorf("GetMembersWithCredit failed.  Expected %v, but found %v", expectedMemberCount, actualMemberCount)
	}
}
