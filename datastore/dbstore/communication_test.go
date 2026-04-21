package dbstore

import (
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
)

func TestGetCommunication_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "name", "subject", "frequency_throttle", "template"}
	mock.ExpectQuery(`from membership\.communication\s+where name = \$1`).
		WithArgs("welcome").
		WillReturnRows(pgxmock.NewRows(cols).AddRow(1, "welcome", "Welcome!", 7, "body"))

	got, err := db.GetCommunication(ctx, "welcome")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "welcome" || got.Subject != "Welcome!" {
		t.Errorf("unexpected communication: %+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetCommunication_Error(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`from membership\.communication`).
		WithArgs("missing").
		WillReturnError(errors.New("no such row"))

	_, err := db.GetCommunication(ctx, "missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetCommunications_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	cols := []string{"id", "name", "subject", "frequency_throttle", "template"}
	mock.ExpectQuery(`Select id, name, subject, frequency_throttle, template\s+from membership\.communication;`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow(1, "welcome", "Welcome!", 7, "body1").
			AddRow(2, "renewal", "Renew", 30, "body2"))

	got := db.GetCommunications(ctx)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

func TestGetMostRecentCommunicationToMember_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	now := time.Now()
	mock.ExpectQuery(`Select created_at\s+from membership\.communication_log`).
		WithArgs("mem-1", 3).
		WillReturnRows(pgxmock.NewRows([]string{"created_at"}).AddRow(now))

	got, err := db.GetMostRecentCommunicationToMember(ctx, "mem-1", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(now) {
		t.Errorf("got %v, want %v", got, now)
	}
}

func TestLogCommunication_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`Insert into membership\.communication_log`).
		WithArgs("mem-1", 7).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err := db.LogCommunication(ctx, 7, "mem-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestLogCommunication_NoRows(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`Insert into membership\.communication_log`).
		WillReturnResult(pgxmock.NewResult("INSERT", 0))

	if err := db.LogCommunication(ctx, 7, "mem-1"); err == nil {
		t.Fatal("expected error when no rows affected")
	}
}

func TestLogCommunication_ExecError(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	mock.ExpectExec(`Insert into membership\.communication_log`).
		WillReturnError(errors.New("duplicate key"))

	if err := db.LogCommunication(ctx, 7, "mem-1"); err == nil {
		t.Fatal("expected error")
	}
}
