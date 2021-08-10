package models

// Communication defines an email communication
type Communication struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Subject           string `json:"subject"`
	FrequencyThrottle int    `json:"frequencyThrottle"`
	Template          string `json:"template"`
}
