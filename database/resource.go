package database

import (
	"context"
	"log"

	"github.com/jackc/pgtype"
)

const getResourceQuery = `SELECT id, description, device_identifier, updated_at FROM membership.resources;`

// Resource -- a resource that can accespt an access control list
type Resource struct {
	ID          uint8            `json:"id"`
	Name        string           `json:"name"`
	Address     string           `json:"address"`
	LastUpdated pgtype.Timestamp `json:"lastUpdated"`
}

// GetResources - gets the status from DB
func (db *Database) GetResources() []Resource {
	rows, err := db.pool.Query(context.Background(), getResourceQuery)
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var resources []Resource

	for rows.Next() {
		var r Resource
		err = rows.Scan(&r.ID, &r.Name, &r.Address, &r.LastUpdated)
		resources = append(resources, r)
	}

	return resources
}
