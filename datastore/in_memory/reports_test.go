package in_memory

import (
	"context"
	"testing"
	"time"

	"github.com/HackRVA/memberserver/models"
)

func TestDemoMemberCounts(t *testing.T) {
	now := time.Date(2026, time.June, 19, 12, 0, 0, 0, time.UTC)
	members := map[string]models.Member{
		"classic-1":  {Level: uint8(models.Classic)},
		"classic-2":  {Level: uint8(models.Classic)},
		"standard-1": {Level: uint8(models.Standard)},
		"premium-1":  {Level: uint8(models.Premium)},
		"credited-1": {Level: uint8(models.Credited)},
		"inactive-1": {Level: uint8(models.Inactive)},
	}

	counts := demoMemberCounts(now, members)
	if len(counts) != 12 {
		t.Fatalf("len = %d, want 12", len(counts))
	}

	latest := counts[len(counts)-1]
	if latest.Month.Month() != time.June || latest.Classic != 2 || latest.Standard != 1 || latest.Premium != 1 || latest.Credited != 1 {
		t.Fatalf("unexpected latest count: %+v", latest)
	}
}

func TestInMemoryReportQueries(t *testing.T) {
	ctx := context.Background()
	jan := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	feb := time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC)
	store := &In_memory{memberCounts: []models.MemberCount{
		{Month: jan, Classic: 5, Standard: 4, Premium: 1, Credited: 2},
		{Month: feb, Classic: 6, Standard: 5, Premium: 2, Credited: 3},
	}}

	got, err := store.GetMemberCountByMonth(ctx, time.Date(2026, time.February, 20, 0, 0, 0, 0, time.UTC))
	if err != nil || got.Classic != 6 {
		t.Fatalf("GetMemberCountByMonth() = %+v, %v", got, err)
	}

	churn, err := store.GetMemberChurn(ctx)
	if err != nil || churn != 3 {
		t.Fatalf("GetMemberChurn() = %d, %v; want 3, nil", churn, err)
	}
}
