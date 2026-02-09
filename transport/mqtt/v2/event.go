package v1

import (
	"encoding/json"
	"fmt"

	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/services/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// {"cmd":"log","type":"access","time":1631240207,"isKnown":"true","access":"Always","username":"Stanley Hash","uid":"f3ec6234","door":"frontdoor"}
type EventLogPayload struct {
	Time     int    `json:"time"`
	Username string `json:"username"`
	RFID     string `json:"uid"`
	Door     string `json:"door"`
}

func (v1 mqttHandler) ReceiveHandler(client mqtt.Client, msg mqtt.Message) {
	var payload models.LogMessage

	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		logger.Errorf("error unmarshalling mqtt payload: %s", err)
		return
	}

	if payload.EventTime == 0 {
		logger.Println("receive: ", string(msg.Payload()))
		return
	}

	if payload.Type == "heartbeat" {
		v1.OnHeartBeatHandler(client, msg)
		return
	}

	v1.OnAccessEventHandler(payload)
}

// OnAccessEvent - post the event to slack. This could also get shoved in the DB eventually
func (v1 mqttHandler) OnAccessEventHandler(payload models.LogMessage) {
	m, err := v1.GetMemberByRFID(payload.RFID)
	if err != nil {
		logger.Debugf("error with access event: %s \n %v %v", err.Error(), m, payload)
		// logger.Errorf("swipe on %s of unknown fob: %s", payload.Door, payload.RFID)
		return
	}

	defer func(m models.Member, p models.LogMessage) {
		go v1.notifier.Send(fmt.Sprintf("name: %s, rfid: %s, door: %s, time: %d", m.Name, p.RFID, p.Door, p.EventTime))
		go func() {
			if err := v1.LogAccessEvent(models.LogMessage{
				Type:      p.Type,
				EventTime: p.EventTime,
				IsKnown:   p.IsKnown,
				Username:  m.Name,
				RFID:      p.RFID,
				Door:      p.Door,
			}); err != nil {
				logger.Errorf("error logging access event %s", err)
			}
		}()
	}(m, payload)
}
