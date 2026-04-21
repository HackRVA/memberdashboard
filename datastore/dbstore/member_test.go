package dbstore

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock"
)

func newTestStore(t *testing.T) (*DatabaseStore, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool: %v", err)
	}
	return &DatabaseStore{ctx: context.Background(), pool: mock}, mock
}

func TestGetMemberCount_Active(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM membership\.members WHERE member_tier_id != 1`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(42))

	got, err := db.GetMemberCount(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 42 {
		t.Errorf("count = %d, want 42", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetMemberCount_Inactive(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM membership\.members WHERE member_tier_id = 1`).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(7))

	got, err := db.GetMemberCount(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 7 {
		t.Errorf("count = %d, want 7", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetMemberCount_QueryError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT COUNT`).WillReturnError(errors.New("boom"))

	got, err := db.GetMemberCount(true)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got != 0 {
		t.Errorf("count = %d, want 0 on error", got)
	}
}
