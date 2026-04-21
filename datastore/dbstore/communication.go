package dbstore

import (
	"context"
	"errors"
	"time"

	"github.com/HackRVA/memberserver/models"

	log "github.com/sirupsen/logrus"
)

var communicationDbMethod CommunicationDatabaseMethod

// GetCommunnications returns all communications from the database
func (db *DatabaseStore) GetCommunications(ctx context.Context) []models.Communication {
	rows, err := db.pool.Query(ctx, communicationDbMethod.getCommunications())
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
func (db *DatabaseStore) GetCommunication(ctx context.Context, name string) (models.Communication, error) {
	var c models.Communication
	err := db.pool.QueryRow(ctx, communicationDbMethod.getCommunication(), name).
		Scan(&c.ID, &c.Name, &c.Subject, &c.FrequencyThrottle, &c.Template)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (db *DatabaseStore) GetMostRecentCommunicationToMember(ctx context.Context, memberId string, commId int) (time.Time, error) {
	var d time.Time
	err := db.pool.QueryRow(ctx, communicationDbMethod.getLastCommunication(), memberId, commId).Scan(&d)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (db *DatabaseStore) LogCommunication(ctx context.Context, communicationId int, memberId string) error {
	commandTag, err := db.pool.Exec(ctx, communicationDbMethod.insertCommunicationLog(), memberId, communicationId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}
