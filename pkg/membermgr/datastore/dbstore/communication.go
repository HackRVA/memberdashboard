package dbstore

import (
	"context"
	"errors"
	"time"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

var communicationDbMethod CommunicationDatabaseMethod

// GetCommunnications returns all communications from the database
func (db *DatabaseStore) GetCommunications() []models.Communication {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	rows, err := dbPool.Query(context.Background(), communicationDbMethod.getCommunications())
	if err != nil {
		log.Errorf("GetCommunications failed: %v", err)
	}

	defer rows.Close()

	var communications []models.Communication

	for rows.Next() {
		var c models.Communication
		err = rows.Scan(&c.ID, &c.Name, &c.Subject, &c.FrequencyThrottle, &c.Template)
		if err == nil {
			communications = append(communications, c)
		}
	}
	return communications
}

// GetCommunnication returns all the requested communication from the database
func (db *DatabaseStore) GetCommunication(name string) (models.Communication, error) {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var c models.Communication
	err = dbPool.QueryRow(context.Background(), communicationDbMethod.getCommunication(), name).
		Scan(&c.ID, &c.Name, &c.Subject, &c.FrequencyThrottle, &c.Template)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (db *DatabaseStore) GetMostRecentCommunicationToMember(memberId string, commId int) (time.Time, error) {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var d time.Time
	err = dbPool.QueryRow(context.Background(), communicationDbMethod.getLastCommunication(), memberId, commId).Scan(&d)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (db *DatabaseStore) LogCommunication(communicationId int, memberId string) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	commandTag, err := dbPool.Exec(context.Background(), communicationDbMethod.insertCommunicationLog(), memberId, communicationId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}
