package dbstore

// ResourceDatabaseMethod -- method container that holds the extension methods to query the resources table
type ResourceDatabaseMethod struct{}

func (resource *ResourceDatabaseMethod) getResource() string {
	return `SELECT id, description, device_identifier, is_default 
	FROM membership.resources
	ORDER BY description;`
}

func (resource *ResourceDatabaseMethod) insertResource() string {
	return `INSERT INTO membership.resources(
		description, device_identifier, is_default)
		VALUES ($1, $2, $3)
		RETURNING *;`
}

func (resource *ResourceDatabaseMethod) updateResource() string {
	return `UPDATE membership.resources
	SET description=$2, device_identifier=$3, is_default=$4
	WHERE id=$1
	RETURNING *;`
}

func (resource *ResourceDatabaseMethod) deleteResource() string {
	return `DELETE FROM membership.resources
	WHERE id = $1;`
}

func (resource *ResourceDatabaseMethod) getResourceByName() string {
	return `SELECT id, description, device_identifier, is_default
	FROM membership.resources
	WHERE description = $1;`
}

func (resource *ResourceDatabaseMethod) getResourceByID() string {
	return `SELECT id, description, device_identifier, is_default
	FROM membership.resources
	WHERE id = $1;`
}

func (resource *ResourceDatabaseMethod) getResourceACLByResourceID() string {
	// only returns the rfid of active members for a specific resourceID
	return `SELECT rfid
	FROM membership.member_resource
	LEFT JOIN membership.members
	ON membership.member_resource.member_id = membership.members.id
	WHERE resource_id = $1
	AND rfid is not NULL
	AND member_tier_id != 1;`
}

func (resource *ResourceDatabaseMethod) getResourceACLByResourceIDQueryWithMemberInfo() string {
	return `SELECT member_id, name, rfid
	FROM membership.member_resource
	LEFT JOIN membership.members
	ON membership.member_resource.member_id = membership.members.id
	WHERE resource_id = $1
	AND rfid is not NULL;`
}

func (resource *ResourceDatabaseMethod) getResourceACLByEmail() string {
	return `SELECT email, device_identifier, description, name, rfid
	FROM membership.member_resource
	LEFT JOIN membership.members
	ON membership.member_resource.member_id = membership.members.id
	LEFT JOIN membership.resources
	ON membership.member_resource.resource_id = membership.resources.id 
	WHERE rfid is not NULL and email = $1;`
}

func (resource *ResourceDatabaseMethod) getMemberResource() string {
	const getMemberResourceQuery = `SELECT id, member_id, resource_id
	FROM membership.member_resource
	WHERE member_id = $1 AND resource_id = $2;`

	return getMemberResourceQuery
}

func (resource *ResourceDatabaseMethod) insertMemberResource() string {
	return `INSERT INTO membership.member_resource(
		member_id, resource_id)
		VALUES ($1, $2)
		RETURNING *;`
}

func (resource *ResourceDatabaseMethod) insertMemberDefaultResource() string {
	return `INSERT INTO membership.member_resource(member_id, resource_id)
	VALUES($1, unnest( ARRAY(SELECT resources.id FROM membership.resources AS resources WHERE resources.is_default IS TRUE)))
	RETURNING *;`
}

func (resource *ResourceDatabaseMethod) removeMemberResource() string {
	return `DELETE FROM membership.member_resource
	WHERE member_id = $1 AND resource_id = $2;`
}
