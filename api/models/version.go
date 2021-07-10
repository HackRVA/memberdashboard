package models

type VersionResponse struct {
	// Commit Hash
	//
	// Example: "ffff"
	Commit string `json:"commit"`
}
