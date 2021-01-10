package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config - values of our config
type Config struct {
	// AccessSecret - secret used for signing jwts
	AccessSecret       string `json:"accessSecret"`
	PaypalClientID     string `json:"paypalClientID"`
	PaypalClientSecret string `json:"paypalClientSecret"`
	PaypalUser         string `json:"paypalUser"`
	PaypalPWD          string `json:"paypalPWD"`
	PaypalSignature    string `json:"paypalSignature"`
	PaypalURL          string `json:"paypalURL"`
	MailgunURL         string `json:"mailgunURL"`
	MailgunKey         string `json:"mailgunKey"`
	MailgunFromAddress string `json:"mailgunFromAddress"`
	MailgunUser        string `json:"mailgunUser"`
	MailgunPassword    string `json:"mailgunPassword"`
}

// Load in the config file to memory
func Load(filepath string) (Config, error) {
	c := Config{}

	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return c, fmt.Errorf("error reading in the config file: %s", err)
	}

	_ = json.Unmarshal([]byte(file), &c)

	return c, err
}
