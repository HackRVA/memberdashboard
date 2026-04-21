package dbstore

import (
	"errors"
	"testing"

	"github.com/HackRVA/memberserver/models"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
)

func TestAddUserToDefaultResources_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	expectMemberByEmailRows(mock, "alice@example.com",
		models.Member{ID: "mem-1", Name: "Alice", Email: "alice@example.com", RFID: "notset", Level: 4},
		[]string{})

	mock.ExpectQuery(`INSERT INTO membership\.member_resource\(member_id, resource_id\)\s+VALUES\(\$1, unnest`).
		WithArgs("mem-1").
		WillReturnRows(pgxmock.NewRows([]string{"id", "member_id", "resource_id"}).
			AddRow("rel-1", "mem-1", "res-1").
			AddRow("rel-2", "mem-1", "res-2"))

	got, err := db.AddUserToDefaultResources("alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].ResourceID != "res-1" {
		t.Errorf("got[0].ResourceID = %q, want res-1", got[0].ResourceID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestAddUserToDefaultResources_MemberNotFound(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("ghost@example.com").
		WillReturnError(pgx.ErrNoRows)

	got, err := db.AddUserToDefaultResources("ghost@example.com")
	if err == nil {
		t.Fatal("expected error")
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestAddMultipleMembersToResource_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	// 1. GetResourceByID
	resCols := []string{"id", "description", "device_identifier", "is_default"}
	mock.ExpectQuery(`FROM membership\.resources\s+WHERE id = \$1`).
		WithArgs("res-1").
		WillReturnRows(pgxmock.NewRows(resCols).AddRow("res-1", "Lasers", "1.1.1.1", false))

	// 2. GetMemberByEmail for alice
	expectMemberByEmailRows(mock, "alice@example.com",
		models.Member{ID: "mem-1", Name: "Alice", Email: "alice@example.com", RFID: "notset", Level: 4},
		[]string{})

	// 3. insertMemberResource
	relCols := []string{"id", "member_id", "resource_id"}
	mock.ExpectQuery(`INSERT INTO membership\.member_resource\(\s+member_id, resource_id\)\s+VALUES \(\$1, \$2\)\s+RETURNING`).
		WithArgs("mem-1", "res-1").
		WillReturnRows(pgxmock.NewRows(relCols).AddRow("rel-1", "mem-1", "res-1"))

	got, err := db.AddMultipleMembersToResource([]string{"alice@example.com"}, "res-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].MemberID != "mem-1" {
		t.Errorf("unexpected result: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestAddMultipleMembersToResource_ResourceNotFound(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.resources\s+WHERE id = \$1`).
		WithArgs("missing").
		WillReturnError(errors.New("not found"))

	_, err := db.AddMultipleMembersToResource([]string{"alice@example.com"}, "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAddMultipleMembersToResource_MemberNotFound(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	resCols := []string{"id", "description", "device_identifier", "is_default"}
	mock.ExpectQuery(`FROM membership\.resources\s+WHERE id = \$1`).
		WithArgs("res-1").
		WillReturnRows(pgxmock.NewRows(resCols).AddRow("res-1", "Lasers", "1.1.1.1", false))

	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("ghost@example.com").
		WillReturnError(pgx.ErrNoRows)

	_, err := db.AddMultipleMembersToResource([]string{"ghost@example.com"}, "res-1")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetResourceACLWithMemberInfo_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"member_id", "name", "rfid"}
	mock.ExpectQuery(`SELECT member_id, name, rfid\s+FROM membership\.member_resource`).
		WithArgs("res-1").
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("mem-1", "Alice", "aaa").
			AddRow("mem-2", "Bob", "bbb"))

	got, err := db.GetResourceACLWithMemberInfo(models.Resource{ID: "res-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Name != "Alice" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestGetMembersAccess_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"email", "device_identifier", "description", "name", "rfid"}
	mock.ExpectQuery(`FROM membership\.member_resource.*WHERE rfid is not NULL and email = \$1`).
		WithArgs("alice@example.com").
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("alice@example.com", "1.1.1.1", "Lasers", "Alice", "aaa"))

	got, err := db.GetMembersAccess(models.Member{Email: "alice@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].ResourceName != "Lasers" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestGetActiveMembersByResource_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"email", "device_identifier", "description", "name", "rfid"}
	mock.ExpectQuery(`FROM membership\.member_resource.*WHERE member_tier_id != 1`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("a@example.com", "1.1.1.1", "Lasers", "Alice", "aaa"))

	got, err := db.GetActiveMembersByResource()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d, want 1", len(got))
	}
}

func TestGetInactiveMembersByResource_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"email", "device_identifier", "description", "name", "rfid"}
	mock.ExpectQuery(`FROM membership\.member_resource.*WHERE member_tier_id = 1`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("inactive@example.com", "1.1.1.1", "Lasers", "Ghost", "xxx"))

	got, err := db.GetInactiveMembersByResource()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len = %d, want 1", len(got))
	}
}
