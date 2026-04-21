package dbstore

import (
	"context"

	config "github.com/HackRVA/memberserver/configs"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// pgxConn is the subset of *pgxpool.Pool that DatabaseStore uses.
// Kept narrow so tests can substitute a mock (e.g. pgxmock).
type pgxConn interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Close()
}

type DatabaseStore struct {
	pool pgxConn
}

func Setup() (*DatabaseStore, error) {
	conf, _ := config.Load()

	pool, err := pgxpool.Connect(context.Background(), conf.DBConnectionString)
	if err != nil {
		return nil, err
	}

	return &DatabaseStore{pool: pool}, nil
}
