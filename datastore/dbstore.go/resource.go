package dbstore

import (
	"errors"
	"fmt"
	"memberserver/api/models"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
)

var resourceDbMethod ResourceDatabaseMethod

var resourceHeartBeatCache map[string]time.Time

// ResourceHeartbeat stores the most recent timestamp that a resource checked in
func ResourceHeartbeat(r models.Resource) {
	if resourceHeartBeatCache == nil {
		resourceHeartBeatCache = make(map[string]time.Time)
	}
	resourceHeartBeatCache[r.Name] = time.Now()
}

// GetLastHeartbeat get the last heart beat
func GetLastHeartbeat(r models.Resource) time.Time {
	return resourceHeartBeatCache[r.Name]
}

// GetResources - gets the status from DB
func (db *DatabaseStore) GetResources() []models.Resource {
	rows, err := db.getConn().Query(db.ctx, resourceDbMethod.getResource())
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var resources []models.Resource

	for rows.Next() {
		var r models.Resource
		_ = rows.Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)

		r.LastHeartBeat = GetLastHeartbeat(r)
		resources = append(resources, r)
	}

	return resources
}

// GetResourceByID - lookup a resource by it's name
func (db *DatabaseStore) GetResourceByID(ID string) (models.Resource, error) {
	var r models.Resource

	err := db.getConn().QueryRow(db.ctx, resourceDbMethod.getResourceByID(), ID).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if err != nil {
		return r, fmt.Errorf("conn.Query failed: %v", err)
	}

	return r, err
}

// GetResourceByName - lookup a resource by it's name
func (db *DatabaseStore) GetResourceByName(resourceName string) (models.Resource, error) {
	var r models.Resource

	err := db.getConn().QueryRow(db.ctx, resourceDbMethod.getResourceByName(), resourceName).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if err != nil {
		return r, fmt.Errorf("getResourceByName failed: %v", err)
	}

	return r, err
}

// RegisterResource - stores a new resource in the db
func (db *DatabaseStore) RegisterResource(name string, address string, isDefault bool) (models.Resource, error) {
	r := &models.Resource{}

	r.Name = name
	r.Address = address
	r.IsDefault = isDefault

	_, err := db.getConn().Exec(db.ctx, resourceDbMethod.insertResource(), r.Name, r.Address, r.IsDefault)
	if err != nil {
		return *r, fmt.Errorf("error inserting resource: %s", err.Error())
	}

	return *r, nil
}

// UpdateResource - updates a resource in the db
func (db *DatabaseStore) UpdateResource(id string, name string, address string, isDefault bool) (*models.Resource, error) {
	r := &models.Resource{}

	// if the resource doesn't already exist let's register it
	if id == "" {
		log.Error("invalid resourseID of 0")
		return r, errors.New("invalid resourseID of 0")
	}

	row := db.getConn().QueryRow(db.ctx, resourceDbMethod.updateResource(), id, name, address, isDefault).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if row == pgx.ErrNoRows {
		log.Printf("no rows affected %s", row.Error())
		return r, errors.New("no rows affected")
	}

	return r, nil
}

// DeleteResource - delete a resource from the db
func (db *DatabaseStore) DeleteResource(id string) error {
	rows, err := db.getConn().Query(db.ctx, resourceDbMethod.deleteResource(), id)
	if err != nil {
		return fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	return nil
}

// AddMultipleMembersToResource grant multiple members access to a resource
func (db *DatabaseStore) AddMultipleMembersToResource(emails []string, resourceID string) ([]models.MemberResourceRelation, error) {

	var membersResource []models.MemberResourceRelation

	resource, err := db.GetResourceByID(resourceID)

	if err != nil {
		return membersResource, err
	}

	for i := 0; i < len(emails); i++ {
		member, err := db.GetMemberByEmail(emails[i])

		if err != nil {
			return membersResource, err
		}

		var memberResource models.MemberResourceRelation
		memberResource.MemberID = member.ID
		memberResource.ResourceID = resource.ID

		row := db.getConn().QueryRow(db.ctx, resourceDbMethod.insertMemberResource(), memberResource.MemberID, memberResource.ResourceID).Scan(&memberResource.ID, &memberResource.MemberID, &memberResource.ResourceID)
		if row == pgx.ErrNoRows {
			return membersResource, errors.New("no rows affected")
		}

		membersResource = append(membersResource, memberResource)

	}

	return membersResource, nil

}

// AddUserToDefaultResources - grants a user access to default resources - untested
func (db *DatabaseStore) AddUserToDefaultResources(email string) ([]models.MemberResourceRelation, error) {
	m, err := db.GetMemberByEmail(email)
	if err != nil {
		return []models.MemberResourceRelation{}, err
	}

	rows, err := db.getConn().Query(db.ctx, resourceDbMethod.insertMemberDefaultResource(), m.ID)
	if err != nil {
		log.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	var memberResources []models.MemberResourceRelation

	for rows.Next() {
		var r models.MemberResourceRelation
		_ = rows.Scan(&r.ID, &r.MemberID, &r.ResourceID)
		memberResources = append(memberResources, r)
	}
	return memberResources, nil
}

// GetMemberResourceRelation retrieves a relation of a member and a resource
func (db *DatabaseStore) GetMemberResourceRelation(m models.Member, r models.Resource) (models.MemberResourceRelation, error) {
	mr := models.MemberResourceRelation{}

	row := db.getConn().QueryRow(db.ctx, resourceDbMethod.getMemberResource(), m.ID, r.ID).Scan(&mr.ID, &mr.MemberID, &mr.ResourceID)
	if row == pgx.ErrNoRows {
		return mr, errors.New("no rows affected")
	}

	return mr, nil
}

// RemoveUserFromResource - removes a users access to a resource
func (db *DatabaseStore) RemoveUserFromResource(email string, resourceID string) error {
	memberResource := models.MemberResourceRelation{}

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

	commandTag, err := db.getConn().Exec(db.ctx, resourceDbMethod.removeMemberResource(), memberResource.MemberID, memberResource.ResourceID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("no row found to delete")
	}

	return nil
}

// GetResourceACL returns a list of members that have access to that Resource
func (db *DatabaseStore) GetResourceACL(r models.Resource) ([]string, error) {
	var accessList []string

	rows, err := db.getConn().Query(db.ctx, resourceDbMethod.getResourceACLByResourceID(), r.ID)
	if err != nil {
		return accessList, fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rfid string
		rows.Scan(&rfid)
		accessList = append(accessList, rfid)
	}

	return accessList, nil
}

// GetResourceACLWithMemberInfo returns a list of members that have access to that Resource
func (db *DatabaseStore) GetResourceACLWithMemberInfo(r models.Resource) ([]models.Member, error) {
	var accessList []models.Member

	rows, err := db.getConn().Query(db.ctx, resourceDbMethod.getResourceACLByResourceIDQueryWithMemberInfo(), r.ID)
	if err != nil {
		return accessList, fmt.Errorf("conn.Query failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var member models.Member

		rows.Scan(&member.ID, &member.Name, &member.RFID)

		accessList = append(accessList, member)
	}

	return accessList, nil
}

// GetMembersAccess returns a list of a specific members access
//   this is used for sending a new rfid assigment to a resource
func (db *DatabaseStore) GetMembersAccess(m models.Member) ([]models.MemberAccess, error) {
	var memberAccess []models.MemberAccess

	rows, err := db.getConn().Query(db.ctx, resourceDbMethod.getResourceACLByEmail(), m.Email)
	if err != nil {
		return memberAccess, fmt.Errorf("error getting members access info: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var resourceUpdate models.MemberAccess

		rows.Scan(&resourceUpdate.Email, &resourceUpdate.ResourceAddress, &resourceUpdate.ResourceName, &resourceUpdate.Name, &resourceUpdate.RFID)

		memberAccess = append(memberAccess, resourceUpdate)
	}

	return memberAccess, nil
}
