package dbstore

import (
	"context"
	"log"
	"memberserver/config"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DatabaseStore - connection pool to communicate to the db
//  in a concurrently safe manner
type DatabaseStore struct {
	pool *pgxpool.Pool
	ctx  context.Context
	mu   sync.Mutex
}

// Setup - sets up connection pool so that we can connect to the db in a
//   concurrently safe manner
func Setup() (*DatabaseStore, error) {
	db := &DatabaseStore{
		ctx:  context.Background(),
		pool: getDBConnection(context.Background()),
	}

	return db, nil
}

// Close - close connection to the db
func (db *DatabaseStore) Release() error {
	ctx := context.Background()
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	return nil
}

func (db *DatabaseStore) getConn() *pgxpool.Pool {

	db.mu.Lock()
	conn, _ := db.pool.Acquire(context.Background())
	if conn.Conn().IsClosed() {
		db.pool = getDBConnection(context.Background())
	}

	defer conn.Release()
	db.mu.Unlock()

	return db.pool
}

func (db *DatabaseStore) PrintStat() {
	stat := db.pool.Stat()
	log.Printf("Pool Stat:\tAquired\tConst\tIdle\tMax\tTotal")
	log.Printf("\t\t%v\t%v\t%v\t%v\t%v", stat.AcquiredConns(), stat.ConstructingConns(), stat.IdleConns(), stat.MaxConns(), stat.TotalConns())
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
	dbPool.Config().ConnConfig.ConnectTimeout = time.Minute * 2
	dbPool.Config().MaxConnIdleTime = time.Minute * 2

	// Get a connection from the pool and check if the database connection is active and working
	db, err := dbPool.Acquire(ctx)

	if err != nil {
		log.Fatalf("failed to get connection on startup: %v\n", err)
	}

	// Add the connection back to the pool
	defer db.Release()

	if err := db.Conn().Ping(ctx); err != nil {
		log.Fatalln(err)
	}

	return dbPool
}
