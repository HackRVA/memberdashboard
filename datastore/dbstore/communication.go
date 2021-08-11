package dbstore

import (
	"context"
	"memberserver/api/models"
	"time"

	log "github.com/sirupsen/logrus"
)

var communicationDbMethod CommunicationDatabaseMethod

// GetCommunnications returns all communications from the database
func (db *DatabaseStore) GetCommunications() []models.Communication {
	rows, err := db.getConn().Query(context.Background(), communicationDbMethod.getCommunications())
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
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
	var c models.Communication
	err := db.getConn().QueryRow(context.Background(), communicationDbMethod.getCommunication(), name).
		Scan(&c.ID, &c.Name, &c.Subject, &c.FrequencyThrottle, &c.Template)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (db *DatabaseStore) GetMostRecentCommunicationToMember(memberId string, commId int) (time.Time, error) {
	var d time.Time
	err := db.getConn().QueryRow(context.Background(), communicationDbMethod.getLastCommunication(), memberId, commId).Scan(&d)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (db *DatabaseStore) LogCommunication(communicationId int, memberId string) error {
	_, err := db.getConn().Exec(context.Background(), communicationDbMethod.insertCommunicationLog(), memberId, communicationId)
	if err != nil {
		return err
	}
	return nil
}
