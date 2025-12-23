package models

type LogMessage struct {
	Type      string `json:"type"`
	EventTime int64  `json:"time"`
	IsKnown   string `json:"isKnown"`
	Username  string `json:"username"`
	RFID      string `json:"uid"`
	Door      string `json:"hostname"`
}
