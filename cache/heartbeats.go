package cache

import (
	"memberserver/database"
	"time"
)

var resources map[string]time.Time

// ResourceHeartBeat stores the most recent timestamp that a resource checked in
func ResourceHeartbeat(r database.Resource) {
	if resources == nil {
		resources = make(map[string]time.Time)
	}
	resources[r.Name] = time.Now()
}

// GetLastHeartbeat
func GetLastHeartbeat(r database.Resource) time.Time {
	return resources[r.Name]
}
