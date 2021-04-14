package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// Config - values of our config
type Config struct {
	ServiceURL          string `json:"serviceURL"`
	ServiceClientID     string `json:"serviceClientID"`
	ServiceClientSecret string `json:"serviceClientSecret"`
	UserName            string `json:"userName"`
	Password            string `json:"password"`
	CredFilePath        string `json:"credFilePath"`
	ResourceName        string `json:"resourceName"`
}

// Load in the config file to memory
//  you can create a config file or pass in Environment variables
//  the config file will take priority
func loadConfig() Config {
	c := Config{}

	// if config file isn't passed in, don't try to look at it
	if len(os.Getenv("KIOSK_CONFIG_FILE")) == 0 {
		return c
	}

	file, err := ioutil.ReadFile(os.Getenv("KIOSK_CONFIG_FILE"))
	if err != nil {
		log.Debugf("error reading in the config file: %s", err)
	}

	_ = json.Unmarshal([]byte(file), &c)

	// if we still don't have an access secret let's generate a random one
	return c
}
