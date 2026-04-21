package dbstore

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HackRVA/memberserver/models"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
)

var resourceDbMethod ResourceDatabaseMethod

var (
	resourceHeartBeatCache = make(map[string]time.Time)
	resourceHeartBeatMu    sync.Mutex
)

// ResourceHeartbeat stores the most recent timestamp that a resource checked in
func ResourceHeartbeat(r models.Resource) {
	resourceHeartBeatMu.Lock()
	resourceHeartBeatCache[r.Name] = time.Now()
	resourceHeartBeatMu.Unlock()
}

// GetLastHeartbeat get the last heart beat
func GetLastHeartbeat(r models.Resource) time.Time {
	resourceHeartBeatMu.Lock()
	t := resourceHeartBeatCache[r.Name]
	resourceHeartBeatMu.Unlock()
	return t
}

// GetResources - gets the status from DB
func (db *DatabaseStore) GetResources(ctx context.Context) []models.Resource {
	rows, err := db.pool.Query(ctx, resourceDbMethod.getResource())
	if err != nil {
		log.Errorf("getResources failed: %v", err)
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
func (db *DatabaseStore) GetResourceByID(ctx context.Context, ID string) (models.Resource, error) {
	var r models.Resource

	err := db.pool.QueryRow(ctx, resourceDbMethod.getResourceByID(), ID).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if err != nil {
		return r, fmt.Errorf("getResourceByID failed: %v", err)
	}

	return r, err
}

// GetResourceByName - lookup a resource by it's name
func (db *DatabaseStore) GetResourceByName(ctx context.Context, resourceName string) (models.Resource, error) {
	var r models.Resource

	err := db.pool.QueryRow(ctx, resourceDbMethod.getResourceByName(), resourceName).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if err != nil {
		return r, fmt.Errorf("getResourceByName failed: %v", err)
	}

	return r, err
}

// RegisterResource - stores a new resource in the db
func (db *DatabaseStore) RegisterResource(ctx context.Context, name string, address string, isDefault bool) (models.Resource, error) {
	r := &models.Resource{}

	r.Name = name
	r.Address = address
	r.IsDefault = isDefault

	commandTag, err := db.pool.Exec(ctx, resourceDbMethod.insertResource(), r.Name, r.Address, r.IsDefault)
	if err != nil {
		return *r, fmt.Errorf("error inserting resource: %s", err.Error())
	}

	if commandTag.RowsAffected() != 1 {
		return *r, errors.New("no row affected")
	}

	return *r, nil
}

// UpdateResource - updates a resource in the db
func (db *DatabaseStore) UpdateResource(ctx context.Context, res models.Resource) (*models.Resource, error) {
	r := &models.Resource{}

	// if the resource doesn't already exist let's register it
	if res.ID == "" {
		log.Error("invalid resourseID of 0")
		return r, errors.New("invalid resourseID of 0")
	}

	row := db.pool.QueryRow(ctx, resourceDbMethod.updateResource(), res.ID, res.Name, res.Address, res.IsDefault).Scan(&r.ID, &r.Name, &r.Address, &r.IsDefault)
	if row == pgx.ErrNoRows {
		log.Printf("no rows affected %s", row.Error())
		return r, errors.New("no rows affected")
	}

	return r, nil
}

// DeleteResource - delete a resource from the db
func (db *DatabaseStore) DeleteResource(ctx context.Context, id string) error {
	rows, err := db.pool.Query(ctx, resourceDbMethod.deleteResource(), id)
	if err != nil {
		return fmt.Errorf("deleteResource failed: %v", err)
	}

	defer rows.Close()

	return nil
}

// AddMultipleMembersToResource grant multiple members access to a resource
func (db *DatabaseStore) AddMultipleMembersToResource(ctx context.Context, emails []string, resourceID string) ([]models.MemberResourceRelation, error) {
	var membersResource []models.MemberResourceRelation

	resource, err := db.GetResourceByID(ctx, resourceID)
	if err != nil {
		return membersResource, err
	}

	for i := 0; i < len(emails); i++ {
		member, err := db.GetMemberByEmail(ctx, emails[i])
		if err != nil {
			return membersResource, err
		}

		var memberResource models.MemberResourceRelation
		memberResource.MemberID = member.ID
		memberResource.ResourceID = resource.ID

		row := db.pool.QueryRow(ctx, resourceDbMethod.insertMemberResource(), memberResource.MemberID, memberResource.ResourceID).Scan(&memberResource.ID, &memberResource.MemberID, &memberResource.ResourceID)
		if row == pgx.ErrNoRows {
			return membersResource, errors.New("no rows affected")
		}

		membersResource = append(membersResource, memberResource)

	}

	return membersResource, nil
}

// AddUserToDefaultResources - grants a user access to default resources - untested
func (db *DatabaseStore) AddUserToDefaultResources(ctx context.Context, email string) ([]models.MemberResourceRelation, error) {
	m, err := db.GetMemberByEmail(ctx, email)
	if err != nil {
		return []models.MemberResourceRelation{}, err
	}

	rows, err := db.pool.Query(ctx, resourceDbMethod.insertMemberDefaultResource(), m.ID)
	if err != nil {
		log.Errorf("addUserToDefaultResources failed: %v", err)
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
func (db *DatabaseStore) GetMemberResourceRelation(ctx context.Context, m models.Member, r models.Resource) (models.MemberResourceRelation, error) {
	mr := models.MemberResourceRelation{}

	row := db.pool.QueryRow(ctx, resourceDbMethod.getMemberResource(), m.ID, r.ID).Scan(&mr.ID, &mr.MemberID, &mr.ResourceID)
	if row == pgx.ErrNoRows {
		return mr, errors.New("no rows affected")
	}

	return mr, nil
}

// RemoveUserFromResource - removes a users access to a resource
func (db *DatabaseStore) RemoveUserFromResource(ctx context.Context, email string, resourceID string) error {
	memberResource := models.MemberResourceRelation{}

	r, err := db.GetResourceByID(ctx, resourceID)
	if err != nil {
		return err
	}

	m, err := db.GetMemberByEmail(ctx, email)
	if err != nil {
		return err
	}

	memberResource, err = db.GetMemberResourceRelation(ctx, m, r)
	if err != nil {
		return err
	}

	commandTag, err := db.pool.Exec(ctx, resourceDbMethod.removeMemberResource(), memberResource.MemberID, memberResource.ResourceID)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no row affected")
	}

	return nil
}

// GetResourceACL returns a list of members that have access to that Resource
func (db *DatabaseStore) GetResourceACL(ctx context.Context, r models.Resource) ([]string, error) {
	var accessList []string

	rows, err := db.pool.Query(ctx, resourceDbMethod.getResourceACLByResourceID(), r.ID)
	if err != nil {
		return accessList, fmt.Errorf("getResourceACL failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var rfid string
		if err := rows.Scan(&rfid); err != nil {
			log.Error(err)
		}
		accessList = append(accessList, rfid)
	}

	return accessList, nil
}

// GetResourceACLWithMemberInfo returns a list of members that have access to that Resource
func (db *DatabaseStore) GetResourceACLWithMemberInfo(ctx context.Context, r models.Resource) ([]models.Member, error) {
	var accessList []models.Member

	rows, err := db.pool.Query(ctx, resourceDbMethod.getResourceACLByResourceIDQueryWithMemberInfo(), r.ID)
	if err != nil {
		return accessList, fmt.Errorf("getResourceACLWithMemberInfo failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var member models.Member

		if err := rows.Scan(&member.ID, &member.Name, &member.RFID); err != nil {
			log.Error(err)
		}

		accessList = append(accessList, member)
	}

	return accessList, nil
}

// GetMembersAccess returns a list of a specific members access
//
//	this is used for sending a new rfid assigment to a resource
func (db *DatabaseStore) GetMembersAccess(ctx context.Context, m models.Member) ([]models.MemberAccess, error) {
	var memberAccess []models.MemberAccess

	rows, err := db.pool.Query(ctx, resourceDbMethod.getResourceACLByEmail(), m.Email)
	if err != nil {
		return memberAccess, fmt.Errorf("error getting members access info: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var resourceUpdate models.MemberAccess

		if err := rows.Scan(&resourceUpdate.Email, &resourceUpdate.ResourceAddress, &resourceUpdate.ResourceName, &resourceUpdate.Name, &resourceUpdate.RFID); err != nil {
			log.Error(err)
		}

		memberAccess = append(memberAccess, resourceUpdate)
	}

	return memberAccess, nil
}

func (db *DatabaseStore) GetInactiveMembersByResource(ctx context.Context) ([]models.MemberAccess, error) {
	var memberAccess []models.MemberAccess

	rows, err := db.pool.Query(ctx, resourceDbMethod.getInactiveMembersResourceACL())
	if err != nil {
		return memberAccess, fmt.Errorf("error getting members access info: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var resourceUpdate models.MemberAccess

		if err := rows.Scan(&resourceUpdate.Email, &resourceUpdate.ResourceAddress, &resourceUpdate.ResourceName, &resourceUpdate.Name, &resourceUpdate.RFID); err != nil {
			log.Error(err)
		}

		memberAccess = append(memberAccess, resourceUpdate)
	}

	return memberAccess, nil
}

func (db *DatabaseStore) GetActiveMembersByResource(ctx context.Context) ([]models.MemberAccess, error) {
	var memberAccess []models.MemberAccess

	rows, err := db.pool.Query(ctx, resourceDbMethod.getActiveMembersResourceACL())
	if err != nil {
		return memberAccess, fmt.Errorf("error getting members access info: %s", err)
	}

	defer rows.Close()

	for rows.Next() {
		var resourceUpdate models.MemberAccess

		if err := rows.Scan(&resourceUpdate.Email, &resourceUpdate.ResourceAddress, &resourceUpdate.ResourceName, &resourceUpdate.Name, &resourceUpdate.RFID); err != nil {
			log.Error(err)
		}

		memberAccess = append(memberAccess, resourceUpdate)
	}

	return memberAccess, nil
}
