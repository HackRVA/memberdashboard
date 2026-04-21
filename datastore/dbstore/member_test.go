package dbstore

import (
	"errors"
	"testing"

	"github.com/HackRVA/memberserver/models"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

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

func TestUpdateMemberByID_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`UPDATE membership\.members SET name=\$1, email=\$2, subscription_id=\$3 WHERE id=\$4`).
		WithArgs("Alice", "a@example.com", "sub-1", "mem-1").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := db.UpdateMemberByID("mem-1", models.Member{Name: "Alice", Email: "a@example.com", SubscriptionID: "sub-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateMemberByID_NoRows(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`UPDATE membership\.members`).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := db.UpdateMemberByID("missing", models.Member{})
	if err == nil {
		t.Fatal("expected error when no rows affected")
	}
}

func TestUpdateMemberByID_ExecError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`UPDATE membership\.members`).WillReturnError(errors.New("db down"))

	err := db.UpdateMemberByID("mem-1", models.Member{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateMemberBySubscriptionID_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`UPDATE membership\.members SET name=\$1, email=\$2 WHERE subscription_id=\$3`).
		WithArgs("Bob", "b@example.com", "sub-9").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := db.UpdateMemberBySubscriptionID("sub-9", models.Member{Name: "Bob", Email: "b@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetMemberByRFID_NoRows(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members\s+WHERE rfid = \$1`).
		WithArgs("deadbeef").
		WillReturnError(pgx.ErrNoRows)

	_, err := db.GetMemberByRFID("deadbeef")
	if err != pgx.ErrNoRows {
		t.Fatalf("err = %v, want pgx.ErrNoRows", err)
	}
}

func TestGetMemberByRFID_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "name", "email", "rfid", "member_tier_id", "resources"}
	mock.ExpectQuery(`FROM membership\.members\s+WHERE rfid = \$1`).
		WithArgs("abc123").
		WillReturnRows(pgxmock.NewRows(cols).AddRow("mem-1", "Alice", "a@example.com", "abc123", uint8(4), []string{}))

	got, err := db.GetMemberByRFID("abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "mem-1" || got.Name != "Alice" {
		t.Errorf("unexpected member: %+v", got)
	}
}

func TestGetTiers_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`SELECT id, description FROM membership\.member_tiers`).
		WillReturnRows(pgxmock.NewRows([]string{"id", "description"}).
			AddRow(uint8(1), "Inactive").
			AddRow(uint8(4), "Standard"))

	tiers := db.GetTiers()
	if len(tiers) != 2 {
		t.Fatalf("len(tiers) = %d, want 2", len(tiers))
	}
	if tiers[1].Name != "Standard" {
		t.Errorf("tiers[1].Name = %q, want Standard", tiers[1].Name)
	}
}

func TestGetMembersWithCredit_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "name", "email", "rfid", "member_tier_id"}
	mock.ExpectQuery(`RIGHT JOIN membership\.member_credit`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("mem-1", "Alice", "a@example.com", "notset", uint8(2)).
			AddRow("mem-2", "Bob", "b@example.com", "notset", uint8(2)))

	got := db.GetMembersWithCredit()
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestUpdateMemberTiers_NoRows(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`UPDATE membership\.members m`).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	db.UpdateMemberTiers()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSetMemberLevel_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "name", "email", "rfid", "member_tier_id", "subscription_id"}
	mock.ExpectQuery(`UPDATE membership\.members\s+SET member_tier_id=\$2\s+WHERE id=\$1`).
		WithArgs("mem-1", models.MemberLevel(2)).
		WillReturnRows(pgxmock.NewRows(cols))

	if err := db.SetMemberLevel("mem-1", models.Credited); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
