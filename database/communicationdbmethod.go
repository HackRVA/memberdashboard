package database

const getCommunications string = `Select id, name, subject, frequency_throttle, template
from membership.communication;`

const getCommunication string = `Select id, name, subject, frequency_throttle, template
from membership.communication
where name = $1;`

const getLastCommunication string = `Select created_at
from membership.communication_log where
member_id = $1
	and communication_id = $2;`

const logCommunication string = `Insert into membership.communication_log
	(member_id, communication_Id)
values
	($1, $2);`
