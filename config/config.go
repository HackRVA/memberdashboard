package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config - values of our config
type Config struct {
	PaypalUser           string `json:"paypalUser"`
	PaypalPWD            string `json:"paypalPWD"`
	PaypalSignature      string `json:"paypalSignature"`
	PaypalURL            string `json:"paypalURL"`
	QBAcceessToken       string `json:"qbAccessToken"`
	QBAcceessTokenSecret string `json:"qbAccessTokenSecret"`
	QBConsumerKey        string `json:"qbConsumerKey"`
	QBConsumerSecret     string `json:"qbConsumerSecret"`
	QBRealmID            string `json:"qbRealmID"`
}

// "PAYPAL_USER"
// "PAYPAL_PWD"
// "PAYPAL_SIGNATURE"
// "PAYPAL_VERSION"
// "PAYPAL_URL"
// 'QB_ACCESS_TOKEN'
// 'QB_ACCESS_TOKEN_SECRET'
// 'QB_CONSUMER_KEY'
// 'QB_CONSUMER_SECRET'
// 'QB_REALM_ID'

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
