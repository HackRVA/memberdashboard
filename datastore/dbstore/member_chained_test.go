package dbstore

import (
	"errors"
	"testing"

	"github.com/HackRVA/memberserver/models"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

// expectMemberQuery stacks a QueryRow response for getMemberByEmail. If
// resourceIDs is non-empty the caller must also expect a GetResourceByID
// query per unique ID.
func expectMemberByEmailRows(mock pgxmock.PgxPoolIface, email string, member models.Member, resourceIDs []string) {
	cols := []string{"id", "name", "email", "rfid", "member_tier_id", "resources"}
	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs(email).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow(member.ID, member.Name, member.Email, member.RFID, member.Level, resourceIDs))
}

func expectResourceByID(mock pgxmock.PgxPoolIface, id string, r models.Resource) {
	cols := []string{"id", "description", "device_identifier", "is_default"}
	mock.ExpectQuery(`FROM membership\.resources\s+WHERE id = \$1`).
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows(cols).AddRow(r.ID, r.Name, r.Address, r.IsDefault))
}

func TestGetMemberByEmail_NoResources(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	expectMemberByEmailRows(mock, "alice@example.com",
		models.Member{ID: "mem-1", Name: "Alice", Email: "alice@example.com", RFID: "notset", Level: 4},
		[]string{})

	got, err := db.GetMemberByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "mem-1" || got.Name != "Alice" {
		t.Errorf("unexpected member: %+v", got)
	}
	if len(got.Resources) != 0 {
		t.Errorf("Resources = %v, want empty", got.Resources)
	}
}

func TestGetMemberByEmail_WithResources(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	expectMemberByEmailRows(mock, "alice@example.com",
		models.Member{ID: "mem-1", Name: "Alice", Email: "alice@example.com", RFID: "notset", Level: 4},
		[]string{"res-1", "res-2"})
	expectResourceByID(mock, "res-1", models.Resource{ID: "res-1", Name: "Lasers"})
	expectResourceByID(mock, "res-2", models.Resource{ID: "res-2", Name: "3D Printer"})

	got, err := db.GetMemberByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Resources) != 2 {
		t.Fatalf("len(Resources) = %d, want 2", len(got.Resources))
	}
	if got.Resources[0].Name != "Lasers" || got.Resources[1].Name != "3D Printer" {
		t.Errorf("unexpected resource names: %+v", got.Resources)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetMemberByEmail_NoRows(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("ghost@example.com").
		WillReturnError(pgx.ErrNoRows)

	_, err := db.GetMemberByEmail("ghost@example.com")
	if err != pgx.ErrNoRows {
		t.Fatalf("err = %v, want pgx.ErrNoRows", err)
	}
}

func TestGetMemberByEmail_QueryError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("alice@example.com").
		WillReturnError(errors.New("conn refused"))

	_, err := db.GetMemberByEmail("alice@example.com")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetMembers_EmptyResources(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "name", "email", "rfid", "member_tier_id", "resources", "subscription_id"}
	mock.ExpectQuery(`FROM membership\.members\s+ORDER BY name`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("mem-1", "Alice", "alice@example.com", "notset", uint8(4), []string{}, "sub-1").
			AddRow("mem-2", "Bob", "bob@example.com", "notset", uint8(4), []string{}, "sub-2"))

	got := db.GetMembers()
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].SubscriptionID != "sub-1" {
		t.Errorf("SubscriptionID = %q, want sub-1", got[0].SubscriptionID)
	}
}

func TestGetMembersPaginated_EmptyResources(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "name", "email", "rfid", "member_tier_id", "resources", "subscription_id"}
	mock.ExpectQuery(`FROM membership\.members\s+WHERE member_tier_id != 1`).
		WithArgs(10, 20). // limit=10, page=2 -> offset=20
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("mem-3", "Carol", "c@example.com", "notset", uint8(4), []string{}, "sub-3"))

	got, err := db.GetMembersPaginated(10, 2, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Carol" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestGetMembersPaginated_QueryError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members`).
		WillReturnError(errors.New("db down"))

	_, err := db.GetMembersPaginated(10, 0, false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAssignRFID_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	// 1. GetMemberByEmail lookup
	expectMemberByEmailRows(mock, "alice@example.com",
		models.Member{ID: "mem-1", Name: "Alice", Email: "alice@example.com", RFID: "notset", Level: 4},
		[]string{})

	// 2. setMemberRFIDTag QueryRow returns the new encoded rfid
	mock.ExpectQuery(`UPDATE membership\.members\s+SET rfid=\$2\s+WHERE email=\$1\s+RETURNING rfid`).
		WithArgs("alice@example.com", encodeRFID("0101436029")).
		WillReturnRows(pgxmock.NewRows([]string{"rfid"}).AddRow(encodeRFID("0101436029")))

	got, err := db.AssignRFID("alice@example.com", "0101436029")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.RFID != encodeRFID("0101436029") {
		t.Errorf("RFID = %q, want %q", got.RFID, encodeRFID("0101436029"))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestAssignRFID_MemberNotFound(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("ghost@example.com").
		WillReturnError(pgx.ErrNoRows)

	_, err := db.AssignRFID("ghost@example.com", "0101436029")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateMember_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	// Lookup returns existing subscription ID; update carries a new one.
	expectMemberByEmailRows(mock, "alice@example.com",
		models.Member{ID: "mem-1", Name: "Alice", Email: "alice@example.com", RFID: "notset", Level: 4},
		[]string{})

	mock.ExpectExec(`UPDATE membership\.members SET name=\$1, subscription_id=\$2 WHERE email=\$3`).
		WithArgs("Alice Updated", "sub-new", "alice@example.com").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := db.UpdateMember(models.Member{
		Name:           "Alice Updated",
		Email:          "alice@example.com",
		SubscriptionID: "sub-new",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateMember_NoRowsAffected(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	expectMemberByEmailRows(mock, "alice@example.com",
		models.Member{ID: "mem-1", Name: "Alice", Email: "alice@example.com", RFID: "notset", Level: 4},
		[]string{})

	mock.ExpectExec(`UPDATE membership\.members`).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := db.UpdateMember(models.Member{Name: "Alice", Email: "alice@example.com"})
	if err == nil {
		t.Fatal("expected error when no rows affected")
	}
}

func TestUpdateMember_LookupFails(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("ghost@example.com").
		WillReturnError(pgx.ErrNoRows)

	err := db.UpdateMember(models.Member{Email: "ghost@example.com"})
	if err == nil {
		t.Fatal("expected error when member lookup fails")
	}
}
