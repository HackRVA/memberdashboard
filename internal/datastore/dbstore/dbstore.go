package dbstore

import (
	"context"
	"memberserver/internal/services/config"
)

type DatabaseStore struct {
	ctx              context.Context
	connectionString string
}

func Setup() (*DatabaseStore, error) {
	conf, _ := config.Load()

	db := &DatabaseStore{
		ctx:              context.Background(),
		connectionString: conf.DBConnectionString,
	}

	return db, nil
}
