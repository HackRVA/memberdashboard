package dbstore

import (
	"context"

	config "github.com/HackRVA/memberserver/configs"
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
