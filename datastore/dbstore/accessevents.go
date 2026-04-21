package dbstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HackRVA/memberserver/models"
)

func (db *DatabaseStore) LogAccessEvent(logMsg models.LogMessage) error {
	timeLayout := "2006-01-02T15:04:05-0700"
	t := time.Unix(logMsg.EventTime, 0)

	commandTag, err := db.pool.Exec(context.Background(), memberDbMethod.insertEvent(), logMsg.Type, t.Format(timeLayout), logMsg.IsKnown, logMsg.Username, logMsg.RFID, logMsg.Door)
	if err != nil {
		return fmt.Errorf("error insterting event to DB: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}
