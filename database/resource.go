package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const getResourceQuery = `SELECT id, description, device_identifier, updated_at FROM membership.resources;`
const getResourceByNameQuery = `SELECT id, description, device_identifier, updated_at
FROM membership.resources
WHERE description = $1;`
const insertResourceQuery = `INSERT INTO membership.resources(
	description, device_identifier, updated_at)
	VALUES ($1, $2, NOW());`
const getMemberResourceQuery = `SELECT id, member_id, resource_id, updated_at
FROM membership.member_resource
WHERE member_id = $1 AND resource_id = $2;`
const insertMemberResourceQuery = `INSERT INTO membership.member_resource(
	member_id, resource_id, updated_at)
	VALUES ($1, $2, NOW());`
const removeMemberResourceQuery = `DELETE FROM membership.member_resource
WHERE id = $1;`

// getAccessListQuery - get a list of rfid tags that belong to an active member
// that have access to a specified resource
const getAccessListQuery = `SELECT rfid
FROM membership.member_resource
INNER JOIN membership.members on (member_resource.member_id = members.id)
WHERE resource_id = $1 AND member_tier_id > 1;`

// Resource -- a resource that can accespt an access control list
type Resource struct {
	ID          uint8            `json:"id"`
	Name        string           `json:"name"`
	Address     string           `json:"address"`
	LastUpdated pgtype.Timestamp `json:"lastUpdated"`
}

// MemberResourceRelation  - a relationship between resources and members
type MemberResourceRelation struct {
	ID          uint8            `json:"id"`
	MemberID    uint8            `json:"resourceID"`
	ResourceID  uint8            `json:"memberID"`
	LastUpdated pgtype.Timestamp `json:"lastUpdated"`
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
		err = rows.Scan(&r.ID, &r.Name, &r.Address, &r.LastUpdated)
		resources = append(resources, r)
	}

	return resources
}

// GetResourceByName - lookup a resource by it's name
func (db *Database) GetResourceByName(resourceName string) (Resource, error) {
	var r Resource
	rows, err := db.pool.Query(context.Background(), getResourceByNameQuery, resourceName)
	if err != nil {
		return r, fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	err = rows.Scan(&r.ID, &r.Name, &r.Address, &r.LastUpdated)
	if err != nil {
		return r, fmt.Errorf("getResourceByName failed: %s", err)
	}

	return r, err
}

// RegisterResource - stores a new resource in the db
func (db *Database) RegisterResource(name string, address string) (*Resource, error) {
	r := &Resource{}

	r.Name = name
	r.Address = address

	row := db.pool.QueryRow(context.Background(), insertResourceQuery, r.Name, r.Address).Scan(&r.ID, &r.Name, &r.Address, &r.LastUpdated)
	if row == pgx.ErrNoRows {
		return r, errors.New("no rows affected")
	}

	return r, nil
}

// AddUserToResource - grants a user access to a resource
func (db *Database) AddUserToResource(email string, resourceName string) (*MemberResourceRelation, error) {
	memberResource := &MemberResourceRelation{}

	r, err := db.GetResourceByName(resourceName)
	if err != nil {
		return memberResource, err
	}

	m, err := db.GetMemberByEmail(email)
	if err != nil {
		return memberResource, err
	}

	memberResource.MemberID = m.ID
	memberResource.ResourceID = r.ID

	row := db.pool.QueryRow(context.Background(), insertMemberResourceQuery, memberResource.MemberID, memberResource.ResourceID).Scan(&memberResource.ID, &memberResource.MemberID, &memberResource.ResourceID, &memberResource.LastUpdated)
	if row == pgx.ErrNoRows {
		return memberResource, errors.New("no rows affected")
	}

	return memberResource, nil
}

// GetMemberResourceRelation retrieves a relation of a member and a resource
func (db *Database) GetMemberResourceRelation(m Member, r Resource) (*MemberResourceRelation, error) {
	mr := &MemberResourceRelation{}

	row := db.pool.QueryRow(context.Background(), getMemberResourceQuery, m.ID, r.ID).Scan(&mr.ID, &mr.MemberID, &mr.ResourceID, &mr.LastUpdated)
	if row == pgx.ErrNoRows {
		return mr, errors.New("no rows affected")
	}

	return mr, nil
}

// RemoveUserFromResource - removes a users access to a resource
func (db *Database) RemoveUserFromResource(email string, resourceName string) (*MemberResourceRelation, error) {
	memberResource := &MemberResourceRelation{}

	r, err := db.GetResourceByName(resourceName)
	if err != nil {
		return memberResource, err
	}

	m, err := db.GetMemberByEmail(email)
	if err != nil {
		return memberResource, err
	}

	memberResource, err = db.GetMemberResourceRelation(m, r)
	if err != nil {
		return memberResource, err
	}

	row := db.pool.QueryRow(context.Background(), removeMemberResourceQuery, memberResource.MemberID, memberResource.ResourceID).Scan(&memberResource.ID, &memberResource.MemberID, &memberResource.ResourceID, &memberResource.LastUpdated)
	if row == pgx.ErrNoRows {
		return memberResource, errors.New("no rows affected")
	}

	return memberResource, nil
}

// GetResourceACL returns a list of members that have access to that Resource
func (db *Database) GetResourceACL(r Resource) ([]uint8, error) {
	var accessList []uint8

	rows, err := db.pool.Query(context.Background(), getResourceQuery)
	if err != nil {
		return accessList, fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rfid uint8
		err = rows.Scan(&rfid)
		accessList = append(accessList, rfid)
	}

	return accessList, nil
}
