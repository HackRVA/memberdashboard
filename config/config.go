package config

import (
	"encoding/json"
	"errors"
	"fmt"
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
}

// Load in the config file to memory
func Load() (Config, error) {
	c := Config{}

	if len(os.Getenv("MEMBER_SERVER_CONFIG_FILE")) == 0 {
		err := errors.New("must set the MEMBER_SERVER_CONFIG_FILE environment variable to point to config file")
		log.Errorf("error loading config: %s", err)
		return c, err
	}

	file, err := ioutil.ReadFile(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))
	if err != nil {
		return c, fmt.Errorf("error reading in the config file: %s", err)
	}

	_ = json.Unmarshal([]byte(file), &c)

	return c, err
}
