package dbstore

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock"
)

func newTestStore(t *testing.T) (*DatabaseStore, pgxmock.PgxPoolIface, context.Context) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool: %v", err)
	}
	return &DatabaseStore{pool: mock}, mock, context.Background()
}
