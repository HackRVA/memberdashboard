package dbstore

func (member *MemberDatabaseMethod) insertEvent() string {
	return `INSERT INTO membership.access_events(
		type, event_time, is_known, username, rfid, door)
		VALUES ($1, $2, $3, $4, $5, $6);;
	`
}
