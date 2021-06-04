package database

// ResourceDatabaseMethod -- method container that holds the extension methods to query the resources table
type ResourceDatabaseMethod struct{}

func (resource *ResourceDatabaseMethod) getResource() string {
	const getResourceQuery = `SELECT id, description, device_identifier, is_default 
	FROM membership.resources
	ORDER BY description;`

	return getResourceQuery
}

func (resource *ResourceDatabaseMethod) insertResource() string {
	const insertResourceQuery = `INSERT INTO membership.resources(
		description, device_identifier, is_default)
		VALUES ($1, $2, $3)
		RETURNING *;`

	return insertResourceQuery
}

func (resource *ResourceDatabaseMethod) updateResource() string {
	const updateResourceQuery = `UPDATE membership.resources
	SET description=$2, device_identifier=$3, is_default=$4
	WHERE id=$1
	RETURNING *;
	`

	return updateResourceQuery
}

func (resource *ResourceDatabaseMethod) deleteResource() string {
	const deleteResourceQuery = `DELETE FROM membership.resources
	WHERE id = $1;`

	return deleteResourceQuery
}

func (resource *ResourceDatabaseMethod) getResourceByName() string {
	const getResourceByNameQuery = `SELECT id, description, device_identifier, is_default
	FROM membership.resources
	WHERE description = $1;`

	return getResourceByNameQuery
}

func (resource *ResourceDatabaseMethod) getResourceByID() string {
	const getResourceByIDQuery = `SELECT id, description, device_identifier, is_default
	FROM membership.resources
	WHERE id = $1;`

	return getResourceByIDQuery
}

func (resource *ResourceDatabaseMethod) getResourceACLByResourceID() string {
	const getResourceACLByResourceIDQuery = `SELECT rfid
	FROM membership.member_resource
	LEFT JOIN membership.members
	ON membership.member_resource.member_id = membership.members.id
	WHERE resource_id = $1
	AND rfid is not NULL;`

	return getResourceACLByResourceIDQuery
}

func (resource *ResourceDatabaseMethod) getResourceACLByResourceIDQueryWithMemberInfo() string {
	const getResourceACLByResourceIDQueryWithMemberInfoQuery = `SELECT member_id, name, rfid
	FROM membership.member_resource
	LEFT JOIN membership.members
	ON membership.member_resource.member_id = membership.members.id
	WHERE resource_id = $1
	AND rfid is not NULL;`

	return getResourceACLByResourceIDQueryWithMemberInfoQuery
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
	const insertMemberResourceQuery = `INSERT INTO membership.member_resource(
		member_id, resource_id)
		VALUES ($1, $2)
		RETURNING *;`

	return insertMemberResourceQuery
}

func (resource *ResourceDatabaseMethod) insertMemberDefaultResource() string {
	const insertMemberDefaultResourceQuery = `INSERT INTO membership.member_resource(member_id, resource_id)
	VALUES($1, unnest( ARRAY(SELECT resources.id FROM membership.resources AS resources WHERE resources.is_default IS TRUE)))
	RETURNING *;`

	return insertMemberDefaultResourceQuery
}

func (resource *ResourceDatabaseMethod) removeMemberResource() string {
	const removeMemberResourceQuery = `DELETE FROM membership.member_resource
	WHERE member_id = $1 AND resource_id = $2;`

	return removeMemberResourceQuery
}

func (resource *ResourceDatabaseMethod) getAccessList() string {
	// getAccessListQuery - get a list of rfid tags that belong to an active member
	// that have access to a specified resource
	const getAccessListQuery = `SELECT rfid
	FROM membership.member_resource
	INNER JOIN membership.members on (member_resource.member_id = members.id)
	WHERE resource_id = $1 AND member_tier_id > 1;`

	return getAccessListQuery
}
