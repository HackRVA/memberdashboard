package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type LogMessage struct {
	Type      string `json:"type"`
	EventTime int64  `json:"time"`
	IsKnown   string `json:"isKnown"`
	Username  string `json:"username"`
	RFID      string `json:"uid"`
	Door      string `json:"door"`
}

func (db *Database) AddLogMsg(logByte []byte) error {
	var logMsg LogMessage

	err := json.Unmarshal(logByte, &logMsg)
	if err != nil {
		return fmt.Errorf("error parsing event: %v", err)
	}

	timeLayout := "2006-01-02T15:04:05-0700"
	t := time.Unix(logMsg.EventTime, 0)
	t.Format(timeLayout)

	_, err = db.getConn().Query(context.Background(), memberDbMethod.insertEvent(), logMsg.Type, t.Format(timeLayout), logMsg.IsKnown, logMsg.Username, logMsg.RFID, logMsg.Door)
	if err != nil {
		return fmt.Errorf("error insterting event to DB: %v", err)
	}

	return nil
}
