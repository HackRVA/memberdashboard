package dbstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func (db *DatabaseStore) LogAccessEvent(logMsg models.LogMessage) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
		return fmt.Errorf("error connecting to DB: %v", err)
	}
	defer dbPool.Close()

	timeLayout := "2006-01-02T15:04:05-0700"
	t := time.Unix(logMsg.EventTime, 0)
	// t.Format(timeLayout)

	commandTag, err := dbPool.Exec(context.Background(), memberDbMethod.insertEvent(), logMsg.Type, t.Format(timeLayout), logMsg.IsKnown, logMsg.Username, logMsg.RFID, logMsg.Door)
	if err != nil {
		return fmt.Errorf("error insterting event to DB: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}
