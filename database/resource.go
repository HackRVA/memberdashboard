package database

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
)

const getResourceQuery = `SELECT id, description, device_identifier, is_default FROM membership.resources;`
const insertResourceQuery = `INSERT INTO membership.resources(
	description, device_identifier, is_default)
	VALUES ($1, $2, $3)
	RETURNING *;`
const updateResourceQuery = `UPDATE membership.resources
SET description=$2, device_identifier=$3, is_default=$4
WHERE id=$1
RETURNING *;
`
const deleteResourceQuery = `DELETE FROM membership.resources
WHERE id = $1;`
const getResourceByNameQuery = `SELECT id, description, device_identifier, is_default
FROM membership.resources
WHERE description = $1;`

const getResourceByIDQuery = `SELECT id, description, device_identifier, is_default
FROM membership.resources
WHERE id = $1;`

const getResourceACLByResourceIDQuery = `SELECT rfid
FROM membership.member_resource
LEFT JOIN membership.members
ON membership.member_resource.member_id = membership.members.id
WHERE resource_id = $1
AND rfid is not NULL;`

const getMemberResourceQuery = `SELECT id, member_id, resource_id
FROM membership.member_resource
WHERE member_id = $1 AND resource_id = $2;`
const insertMemberResourceQuery = `INSERT INTO membership.member_resource(
	member_id, resource_id)
	VALUES ($1, $2)
	RETURNING *;`
const insertMemberDefaultResourceQuery = `INSERT INTO membership.member_resource(member_id, resource_id)
SELECT $1, resources.id FROM membership.resources AS resources WHERE resources.is_default IS TRUE
RETURNING *;`
const removeMemberResourceQuery = `DELETE FROM membership.member_resource
WHERE member_id = $1 AND resource_id = $2;`

// getAccessListQuery - get a list of rfid tags that belong to an active member
// that have access to a specified resource
const getAccessListQuery = `SELECT rfid
FROM membership.member_resource
INNER JOIN membership.members on (member_resource.member_id = members.id)
WHERE resource_id = $1 AND member_tier_id > 1;`

// Resource a resource that can accespt an access control list
type Resource struct {
	// UniqueID of the Resource
	// required: true
	// example: 0
	ID string `json:"id"`
	// Name of the Resource
	// required: true
	// example: name
	Name string `json:"name"`
	// Address of the Resource. i.e. where it can be found on the network
	// required: true
	// example: address
	Address string `json:"address"`
    // Default state of the Resource
    // required: true
    // example: true
    IsDefault bool `json:"is_default"`
}

// ResourceDeleteRequest - request for deleting a resource
type ResourceDeleteRequest struct {
	// UniqueID of the Resource
	// required: true
	// example: ""
	ID string `json:"id"`
}

// Resource a resource that can accespt an access control list
type ResourceRequest struct {
	// UniqueID of the Resource
	// required: true
	// example: 0
	ID string `json:"id"`
	// Name of the Resource
	// required: true
	// example: name
	Name string `json:"name"`
	// Address of the Resource. i.e. where it can be found on the network
	// required: true
	// example: address
	Address string `json:"address"`
    // Default state of the Resource
    // required: true
    // example: true
    IsDefault bool `json:"is_default"`
}

// Resource a resource that can accespt an access control list
type RegisterResourceRequest struct {
	// Name of the Resource
	// required: true
	// example: name
	Name string `json:"name"`
	// Address of the Resource. i.e. where it can be found on the network
	// required: true
	// example: address
	Address string `json:"address"`
    // Default state of the Resource
    // required: false
    // example: true
    IsDefault bool `json:"is_default"`
}

// MemberResourceRelation  - a relationship between resources and members
type MemberResourceRelation struct {
	ID         string `json:"id"`
	MemberID   string `json:"memberID"`
	ResourceID string `json:"resourceID"`
}

// GetResources - gets the status from DB
func (db *Database) GetResources() []Resource {
	rows, err := db.pool.Query(context.Background(), getResourceQuery)
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var resources []Resource

	for rows.Next() {
		var r Resource
		err = rows.Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
		resources = append(resources, r)
	}

	return resources
}

// GetResourceByID - lookup a resource by it's name
func (db *Database) GetResourceByID(ID string) (Resource, error) {
	var r Resource
	err := db.pool.QueryRow(context.Background(), getResourceByIDQuery, ID).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if err != nil {
		return r, fmt.Errorf("conn.Query failed: %v", err)
	}

	return r, err
}

// GetResourceByName - lookup a resource by it's name
func (db *Database) GetResourceByName(resourceName string) (Resource, error) {
	var r Resource
    err := db.pool.QueryRow(context.Background(), getResourceByNameQuery, resourceName).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if err != nil {
		return r, fmt.Errorf("getResourceByName failed: %v", err)
	}

	return r, err
}

// RegisterResource - stores a new resource in the db
func (db *Database) RegisterResource(name string, address string, is_default bool) (*Resource, error) {
	r := &Resource{}

	r.Name = name
	r.Address = address
    r.IsDefault = is_default

	_, err := db.pool.Exec(context.Background(), insertResourceQuery, r.Name, r.Address, r.IsDefault)
	if err != nil {
		return r, fmt.Errorf("error inserting resource: %s", err.Error())
	}

	return r, nil
}

// UpdateResource - updates a resource in the db
func (db *Database) UpdateResource(id string, name string, address string, is_default bool) (*Resource, error) {
	r := &Resource{}

	// if the resource doesn't already exist let's register it
	if id == "" {
		log.Error("invalid resourseID of 0")
		return r, errors.New("invalid resourseID of 0")
	}

    row := db.pool.QueryRow(context.Background(), updateResourceQuery, id, name, address, is_default).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if row == pgx.ErrNoRows {
		log.Printf("no rows affected %s", row.Error())
		return r, errors.New("no rows affected")
	}

	return r, nil
}

// DeleteResource - delete a resource from the db
func (db *Database) DeleteResource(id string) error {
	rows, err := db.pool.Query(context.Background(), deleteResourceQuery, id)
	if err != nil {
		return fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	return nil
}

// AddUserToResource - grants a user access to a resource
func (db *Database) AddUserToResource(email string, resourceID string) (MemberResourceRelation, error) {
	memberResource := MemberResourceRelation{}

	r, err := db.GetResourceByID(resourceID)
	if err != nil {
		return memberResource, err
	}

	m, err := db.GetMemberByEmail(email)
	if err != nil {
		return memberResource, err
	}

	memberResource.MemberID = m.ID
	memberResource.ResourceID = r.ID

	row := db.pool.QueryRow(context.Background(), insertMemberResourceQuery, memberResource.MemberID, memberResource.ResourceID).Scan(&memberResource.ID, &memberResource.MemberID, &memberResource.ResourceID)
	if row == pgx.ErrNoRows {
		return memberResource, errors.New("no rows affected")
	}

	return memberResource, nil
}

// AddUserToDefaultResources - grants a user access to default resources - untested
func (db *Database) AddUserToDefaultResources(email string) ([]MemberResourceRelation, error) {

	m, err := db.GetMemberByEmail(email)
	if err != nil {
        //is this wrong?
		return [], err
	}

	rows, err := db.pool.Query(context.Background(), insertMemberDefaultResourceQuery, memberResource.MemberID)
	if err != nil {
		log.Fatalf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var memberResources := []MemberResourceRelation

	for rows.Next() {
		var r MemberResourceRelation
		err = rows.Scan(&r.ID, &r.MemberID, &r.ResourceID)
		memberResources = append(memberResources, r)
	}
	return memberResources, nil
}

// GetMemberResourceRelation retrieves a relation of a member and a resource
func (db *Database) GetMemberResourceRelation(m Member, r Resource) (MemberResourceRelation, error) {
	mr := MemberResourceRelation{}

	row := db.pool.QueryRow(context.Background(), getMemberResourceQuery, m.ID, r.ID).Scan(&mr.ID, &mr.MemberID, &mr.ResourceID)
	if row == pgx.ErrNoRows {
		return mr, errors.New("no rows affected")
	}

	return mr, nil
}

// RemoveUserFromResource - removes a users access to a resource
func (db *Database) RemoveUserFromResource(email string, resourceID string) error {
	memberResource := MemberResourceRelation{}

	r, err := db.GetResourceByID(resourceID)
	if err != nil {
		return err
	}

	m, err := db.GetMemberByEmail(email)
	if err != nil {
		return err
	}

	memberResource, err = db.GetMemberResourceRelation(m, r)
	if err != nil {
		return err
	}

	commandTag, err := db.pool.Exec(context.Background(), removeMemberResourceQuery, memberResource.MemberID, memberResource.ResourceID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}

	return nil
}

// GetResourceACL returns a list of members that have access to that Resource
func (db *Database) GetResourceACL(r Resource) ([]string, error) {
	var accessList []string

	rows, err := db.pool.Query(context.Background(), getResourceACLByResourceIDQuery, r.ID)
	if err != nil {
		return accessList, fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rfid string
		err = rows.Scan(&rfid)
		accessList = append(accessList, rfid)
	}

	return accessList, nil
}
