package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config - values of our config
type Config struct {
	TestValue     string `json:"testValue"`
	SomethingElse string `json:"somethingElse"`
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
