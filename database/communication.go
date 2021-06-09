package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

const getCommunications string = "Select id, name, subject, frequency_throttle, template from membership.communication;"
const getCommunication string = "Select id, name, subject, frequency_throttle, template from membership.communication where name = $1;"
const getLastCommunication string = "Select created_at from membership.communication_log where member_id = $1 and communication_id = $2;"
const logCommunication string = "Insert into membership.communication_log (member_id, communication_Id) values ($1, $2);"

// Communication defines an email communication
type Communication struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Subject           string `json:"subject"`
	FrequencyThrottle int    `json:"frequency_throttle"`
	Template          string `json:"template"`
}

// GetCommunnications returns all communications from the database
func (db *Database) GetCommunications() []Communication {
	rows, err := db.getConn().Query(context.Background(), getCommunications)
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var communications []Communication

	for rows.Next() {
		var c Communication
		err = rows.Scan(&c.ID, &c.Name, &c.Subject, &c.FrequencyThrottle, &c.Template)
		communications = append(communications, c)
	}
	return communications
}

// GetCommunnication returns all the requested communication from the database
func (db *Database) GetCommunication(name string) (Communication, error) {
	var c Communication
	err := db.getConn().QueryRow(context.Background(), getCommunication, name).
		Scan(&c.ID, &c.Name, &c.Subject, &c.FrequencyThrottle, &c.Template)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (db *Database) GetMostRecentCommunicationToMember(memberId string, commId int) (time.Time, error) {
	var d time.Time
	err := db.getConn().QueryRow(context.Background(), getCommunication, memberId, commId).Scan(&d)
	if err != nil {
		return d, err
	}
	return d, nil
}

func (db *Database) LogCommunication(communicationId int, memberId string) error {
	_, err := db.getConn().Exec(context.Background(), logCommunication, communicationId, memberId)
	if err != nil {
		return err
	}
	return nil
}
