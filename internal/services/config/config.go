package config

import (
	"encoding/json"
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
	MQTTBrokerAddress  string `json:"mqttBrokerAddress"`
	DBConnectionString string `json:"dbConnectionString"`
	// EnableInfoEmails sends notification by email to info@hackrva.org
	EnableInfoEmails bool `json:"enableInfoEmails"`
	// EnableNotificationEmailsToMembers sends notification to membership
	EnableNotificationEmailsToMembers bool `json:"enableNotificationEmailsToMembers"`
	// EmailOverrideAddress config can provide an email address to send to instead of
	//   the predefined addresses
	EmailOverrideAddress string `json:"emailOverrideAddress"`
	SlackToken           string `json:"slackToken"`
	AdminEmail           string `json:"adminEmail"`
	AlwaysAdmin          string `json:"alwaysAdmin"`
}

// Get gets the config and ignores errors
func Get() Config {
	c, _ := Load()
	return c
}

// Load in the config file to memory
//  you can create a config file or pass in Environment variables
//  the config file will take priority
func Load() (Config, error) {
	c := Config{}

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
	c.MQTTBrokerAddress = os.Getenv("MQTT_BROKER_ADDRESS")
	c.EmailOverrideAddress = os.Getenv("EMAIL_OVERRIDE_ADDRESS")
	c.EnableInfoEmails = false
	c.EnableNotificationEmailsToMembers = false
	c.SlackToken = os.Getenv("SLACK_TOKEN")
	c.AdminEmail = getEnvOrDefault("ADMIN_EMAIL", "info@hackrva.org")
	c.AlwaysAdmin = getEnvOrDefault("ALWAYS_ADMIN", "false")

	if len(os.Getenv("ENABLE_INFO_EMAILS")) > 0 {
		c.EnableInfoEmails = true
	}
	if len(os.Getenv("ENABLE_MEMBER_EMAILS")) > 0 {
		c.EnableNotificationEmailsToMembers = true
	}

	c.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")

	if os.Getenv("DATABASE_URL") != "" {
		c.DBConnectionString = os.Getenv("DATABASE_URL")
	}

	// if config file isn't passed in, don't try to look at it
	if len(os.Getenv("MEMBER_SERVER_CONFIG_FILE")) == 0 {
		return c, nil
	}

	file, err := ioutil.ReadFile(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))
	if err != nil {
		log.Debugf("error reading in the config file: %s", err)
	}

	if err := json.Unmarshal([]byte(file), &c); err != nil {
		return Config{}, err
	}

	// if we still don't have an access secret let's generate a random one
	return c, nil
}

func getEnvOrDefault(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}
