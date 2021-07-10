package models

type VersionResponse struct {
	Major  string `json:"major"`
	Minor  string `json:"minor"`
	Hotfix string `json:"hotfix"`
	Build  string `json:"build"`
}
