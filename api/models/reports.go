package models

import "time"

type MemberCount struct {
	Month    time.Time `json:"month"`
	Classic  int       `json:"classic"`
	Standard int       `json:"standard"`
	Premium  int       `json:"premium"`
	Credited int       `json:"credited"`
}

type AccessStats struct {
	Date         time.Time `json:"date"`
	AccessCount  int       `json:"accessCount"`
	ResourceName string    `json:"resourceName"`
}
