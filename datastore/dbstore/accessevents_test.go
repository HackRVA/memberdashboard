package dbstore

import (
	"errors"
	"testing"

	"github.com/HackRVA/memberserver/models"

	"github.com/pashagolub/pgxmock"
)

func TestLogAccessEvent_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO membership\.access_events`).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := db.LogAccessEvent(models.LogMessage{
		Type:      "access",
		EventTime: 1700000000,
		IsKnown:   "true",
		Username:  "alice",
		RFID:      "deadbeef",
		Door:      "front",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestLogAccessEvent_NoRows(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO membership\.access_events`).
		WillReturnResult(pgxmock.NewResult("INSERT", 0))

	err := db.LogAccessEvent(models.LogMessage{Type: "access"})
	if err == nil {
		t.Fatal("expected error when no rows affected")
	}
}

func TestLogAccessEvent_ExecError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO membership\.access_events`).
		WillReturnError(errors.New("constraint violation"))

	err := db.LogAccessEvent(models.LogMessage{Type: "access"})
	if err == nil {
		t.Fatal("expected error")
	}
}
