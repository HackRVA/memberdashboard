package dbstore

import (
	"errors"
	"testing"

	"github.com/HackRVA/memberserver/models"

	"github.com/pashagolub/pgxmock"
)

func TestGetResourceByID_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "description", "device_identifier", "is_default"}
	mock.ExpectQuery(`FROM membership\.resources\s+WHERE id = \$1`).
		WithArgs("res-1").
		WillReturnRows(pgxmock.NewRows(cols).AddRow("res-1", "Lasers", "192.168.1.2", true))

	got, err := db.GetResourceByID("res-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Lasers" || !got.IsDefault {
		t.Errorf("unexpected resource: %+v", got)
	}
}

func TestGetResourceByID_Error(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.resources\s+WHERE id = \$1`).
		WithArgs("missing").
		WillReturnError(errors.New("no rows"))

	_, err := db.GetResourceByID("missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetResourceByName_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "description", "device_identifier", "is_default"}
	mock.ExpectQuery(`FROM membership\.resources\s+WHERE description = \$1`).
		WithArgs("Lasers").
		WillReturnRows(pgxmock.NewRows(cols).AddRow("res-1", "Lasers", "192.168.1.2", false))

	got, err := db.GetResourceByName("Lasers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Address != "192.168.1.2" {
		t.Errorf("Address = %q, want 192.168.1.2", got.Address)
	}
}

func TestRegisterResource_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO membership\.resources`).
		WithArgs("Lasers", "1.2.3.4", true).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	got, err := db.RegisterResource("Lasers", "1.2.3.4", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Lasers" || got.Address != "1.2.3.4" || !got.IsDefault {
		t.Errorf("unexpected resource: %+v", got)
	}
}

func TestRegisterResource_NoRows(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO membership\.resources`).
		WillReturnResult(pgxmock.NewResult("INSERT", 0))

	_, err := db.RegisterResource("Lasers", "1.2.3.4", false)
	if err == nil {
		t.Fatal("expected error when no rows affected")
	}
}

func TestRegisterResource_ExecError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`INSERT INTO membership\.resources`).
		WillReturnError(errors.New("unique violation"))

	_, err := db.RegisterResource("Lasers", "1.2.3.4", false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateResource_EmptyID(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	_, err := db.UpdateResource(models.Resource{ID: ""})
	if err == nil {
		t.Fatal("expected error for empty resource ID")
	}
}

func TestDeleteResource_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`DELETE FROM membership\.resources`).
		WithArgs("res-1").
		WillReturnRows(pgxmock.NewRows([]string{}))

	if err := db.DeleteResource("res-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteResource_QueryError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`DELETE FROM membership\.resources`).
		WillReturnError(errors.New("fk violation"))

	if err := db.DeleteResource("res-1"); err == nil {
		t.Fatal("expected error")
	}
}

func TestGetResources_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "description", "device_identifier", "is_default"}
	mock.ExpectQuery(`FROM membership\.resources\s+ORDER BY description`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow("res-1", "Lasers", "1.1.1.1", true).
			AddRow("res-2", "3D Printer", "1.1.1.2", false))

	got := db.GetResources()
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestGetResourceACL_Success(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.member_resource`).
		WithArgs("res-1").
		WillReturnRows(pgxmock.NewRows([]string{"rfid"}).
			AddRow("aaa").
			AddRow("bbb"))

	got, err := db.GetResourceACL(models.Resource{ID: "res-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "aaa" {
		t.Errorf("unexpected ACL: %+v", got)
	}
}

func TestGetResourceACL_QueryError(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.member_resource`).
		WillReturnError(errors.New("db down"))

	_, err := db.GetResourceACL(models.Resource{ID: "res-1"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRemoveUserFromResource_DeleteSucceeds(t *testing.T) {
	db, mock := newTestStore(t)
	defer mock.Close()

	resCols := []string{"id", "description", "device_identifier", "is_default"}
	mock.ExpectQuery(`FROM membership\.resources\s+WHERE id = \$1`).
		WithArgs("res-1").
		WillReturnRows(pgxmock.NewRows(resCols).AddRow("res-1", "Lasers", "1.1.1.1", false))

	memberCols := []string{"id", "name", "email", "rfid", "member_tier_id", "resources"}
	mock.ExpectQuery(`FROM membership\.members\s+WHERE LOWER\(email\) = LOWER\(\$1\)`).
		WithArgs("alice@example.com").
		WillReturnRows(pgxmock.NewRows(memberCols).
			AddRow("mem-1", "Alice", "alice@example.com", "rfid", uint8(4), []string{}))

	relCols := []string{"id", "member_id", "resource_id"}
	mock.ExpectQuery(`FROM membership\.member_resource\s+WHERE member_id = \$1 AND resource_id = \$2`).
		WithArgs("mem-1", "res-1").
		WillReturnRows(pgxmock.NewRows(relCols).AddRow("rel-1", "mem-1", "res-1"))

	mock.ExpectExec(`DELETE FROM membership\.member_resource\s+WHERE member_id = \$1 AND resource_id = \$2`).
		WithArgs("mem-1", "res-1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err := db.RemoveUserFromResource("alice@example.com", "res-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
