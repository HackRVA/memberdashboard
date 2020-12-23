package database

import (
	"context"
	"log"
)

const getStatusQuery = `SELECT id, description FROM membership.member_statuses;`

// GetStatuses - gets the status from DB
func (db *Database) GetStatuses() []string {
	rows, err := db.pool.Query(context.Background(), getStatusQuery)
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var statusList []string

	for rows.Next() {
		var id int32
		var status string
		err = rows.Scan(&id, &status)
		if err != nil {
			log.Fatal(err)
		}

		statusList = append(statusList, status)
	}

	return statusList
}
