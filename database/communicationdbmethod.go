package database

type CommunicationDatabaseMethod struct{}

func (CommunicationDatabaseMethod) getCommunications() string {
	return `Select id, name, subject, frequency_throttle, template
from membership.communication;`
}

func (CommunicationDatabaseMethod) getCommunication() string {
	return `Select id, name, subject, frequency_throttle, template
from membership.communication
where name = $1;`
}

func (CommunicationDatabaseMethod) getLastCommunication() string {
	return `Select created_at
from membership.communication_log where
member_id = $1
	and communication_id = $2;`
}

func (CommunicationDatabaseMethod) insertCommunicationLog() string {
	return `Insert into membership.communication_log
	(member_id, communication_Id)
values
	($1, $2);`
}
