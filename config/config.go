package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// Config - values of our config
type Config struct {
	// AccessSecret - secret used for signing jwts
	AccessSecret       string `json:"accessSecret"`
	PaypalClientID     string `json:"paypalClientID"`
	PaypalClientSecret string `json:"paypalClientSecret"`
	PaypalURL          string `json:"paypalURL"`
	MailgunURL         string `json:"mailgunURL"`
	MailgunKey         string `json:"mailgunKey"`
	MailgunFromAddress string `json:"mailgunFromAddress"`
	MailgunUser        string `json:"mailgunUser"`
	MailgunPassword    string `json:"mailgunPassword"`
	SlackAccessEvents  string `json:"slackhookAccessEvents"`
	MQTTUsername       string `json:"mqttUsername"`
	MQTTPassword       string `json:"mqttPassword"`
	DBConnectionString string `json:"dbConnectionString"`
}

// Load in the config file to memory
//  you can create a config file or pass in Environment variables
//  the config file will take priority
func Load() (Config, error) {
	c := Config{}

	if len(os.Getenv("MEMBER_SERVER_CONFIG_FILE")) == 0 {
		err := errors.New("must set the MEMBER_SERVER_CONFIG_FILE environment variable to point to config file")
		log.Debugf("error loading config: %s  Attempting to use environment variables", err)
	}

	file, err := ioutil.ReadFile(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))
	if err != nil {
		log.Debugf("error reading in the config file: %s", err)
	}

	c.AccessSecret = os.Getenv("ACCESS_SECRET")
	c.PaypalClientID = os.Getenv("PAYPAL_CLIENT_ID")
	c.PaypalClientSecret = os.Getenv("PAYPAL_CLIENT_SECRET")
	c.PaypalURL = os.Getenv("PAYPAL_API_URL")
	c.MailgunURL = os.Getenv("MAILGUN_API_URL")
	c.MailgunKey = os.Getenv("MAILGUN_KEY")
	c.MailgunFromAddress = os.Getenv("MAILGUN_FROM_ADDRESS")
	c.MailgunUser = os.Getenv("MAILGUN_USER")
	c.MailgunPassword = os.Getenv("MAILGUN_PASSWORD")
	c.SlackAccessEvents = os.Getenv("SLACK_ACCESS_EVENTS_HOOK")
	c.MQTTUsername = os.Getenv("MQTT_USERNAME")
	c.MQTTPassword = os.Getenv("MQTT_PASSWORD")
	c.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")

	_ = json.Unmarshal([]byte(file), &c)

	// if we still don't have an access secret let's generate a random one
	return c, nil
}
