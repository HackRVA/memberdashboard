package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database - connection pool to communicate to the db
//  in a concurrently safe manner
type Database struct {
	pool *pgxpool.Pool
}

func postgreSQLDatabase() (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to DB!")
	return conn, err
}

// Setup - sets up connection pool so that we can connect to the db in a
//   concurrently safe manner
func Setup() (*Database, error) {
	var err error

	db := &Database{}

	connection, err := postgreSQLDatabase()
	if err != nil {
		return db, err
	}

	db.pool = connection

	return db, nil
}

// Close - close connection to the db
func (db *Database) Release() error {
	ctx := context.Background()
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	return nil
}
