package dbstore

import (
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
)

func TestGetMemberCounts_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	jan := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	feb := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	cols := []string{"month", "classic", "standard", "premium", "credited"}
	mock.ExpectQuery(`FROM membership\.member_counts\s+ORDER BY month`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow(jan, 1, 2, 3, 4).
			AddRow(feb, 5, 6, 7, 8))

	got, err := db.GetMemberCounts(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[1].Premium != 7 {
		t.Errorf("got[1].Premium = %d, want 7", got[1].Premium)
	}
}

func TestGetMemberCounts_QueryError(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.member_counts`).
		WillReturnError(errors.New("db down"))

	_, err := db.GetMemberCounts(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetMemberCountByMonth_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	month := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	cols := []string{"classic", "standard", "premium", "credited"}
	mock.ExpectQuery(`FROM membership\.member_counts\s+WHERE month`).
		WithArgs(month).
		WillReturnRows(pgxmock.NewRows(cols).AddRow(10, 20, 30, 40))

	got, err := db.GetMemberCountByMonth(ctx, month)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Classic != 10 || got.Premium != 30 {
		t.Errorf("unexpected counts: %+v", got)
	}
}

func TestGetMemberChurn_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	jan := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	feb := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	cols := []string{"month", "member_count"}
	mock.ExpectQuery(`FROM membership\.member_counts\s+ORDER BY month DESC`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow(feb, 95).
			AddRow(jan, 100))

	churn, err := db.GetMemberChurn(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if churn != -5 {
		t.Errorf("churn = %d, want -5", churn)
	}
}

func TestGetMemberChurn_QueryError(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	mock.ExpectQuery(`FROM membership\.member_counts`).
		WillReturnError(errors.New("db down"))

	_, err := db.GetMemberChurn(ctx)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAccessStats_Success(t *testing.T) {
	db, mock, ctx := newTestStore(t)
	defer mock.Close()

	day := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	cols := []string{"day", "resource", "count"}
	mock.ExpectQuery(`FROM membership\.access_events`).
		WillReturnRows(pgxmock.NewRows(cols).
			AddRow(day, "front", 12))

	got, err := db.GetAccessStats(ctx, day, "front")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].AccessCount != 12 {
		t.Errorf("unexpected stats: %+v", got)
	}
}
