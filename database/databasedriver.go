package database

import (
	"context"

	log "github.com/sirupsen/logrus"

	"memberserver/config"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database - connection pool to communicate to the db
//  in a concurrently safe manner
type Database struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

// Setup - sets up connection pool so that we can connect to the db in a
//   concurrently safe manner
func Setup() (*Database, error) {
	db := &Database{
		ctx:  context.Background(),
		pool: getDBConnection(context.Background()),
	}

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

func (db *Database) getConn() *pgxpool.Pool {
	conn, _ := db.pool.Acquire(context.Background())
	if conn.Conn().IsClosed() {
		db.pool = getDBConnection(context.Background())
	}

	defer conn.Release()
	return db.pool
}

func getDBConnection(ctx context.Context) *pgxpool.Pool {
	// Retrieve the database host address
	conf, _ := config.Load()

	dbPool, err := pgxpool.Connect(ctx, conf.DBConnectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}

	// Try connecting to the database a few times before giving up
	// Retry to connect for a while
	for i := 1; i < 8 && err != nil; i++ {
		// Sleep a bit before trying again
		time.Sleep(time.Duration(i*i) * time.Second)

		log.Printf("trying to connect to the db server (attempt %d)...\n", i)

		dbPool, err = pgxpool.Connect(ctx, conf.DBConnectionString)
		if err != nil {
			log.Printf("got error: %v\n", err)
		}
	}

	// Stop execution if the database was not initialized
	if dbPool == nil {
		log.Fatalln("could not connect to the database")
	}

	dbPool.Config().ConnConfig.PreferSimpleProtocol = true

	// Get a connection from the pool and check if the database connection is active and working
	db, err := dbPool.Acquire(ctx)
	if err != nil {
		log.Fatalf("failed to get connection on startup: %v\n", err)
	}
	if err := db.Conn().Ping(ctx); err != nil {
		log.Fatalln(err)
	}

	// Add the connection back to the pool
	db.Release()

	return dbPool
}
