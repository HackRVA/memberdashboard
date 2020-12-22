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

var (
	// DB instance of the database connection
	DB *Database
)

func postgreSQLDatabase() (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected!")
	return conn, err
}

// newDB - sets up connection pool so that we can connect to the db in a
//   concurrently safe manner
func newDB() (*Database, error) {
	var err error

	db := &Database{}

	connection, err := postgreSQLDatabase()
	if err != nil {
		return db, err
	}

	db.pool = connection

	return db, nil
}

// Setup -- initialize the db connection
func Setup() error {
	var err error
	DB, err = newDB()

	if err != nil {
		return fmt.Errorf("Error initializing db connection: %s", err.Error())
	}

	return err
}
