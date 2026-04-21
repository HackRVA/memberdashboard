package dbstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/HackRVA/memberserver/services/logger"

	"github.com/HackRVA/memberserver/models"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

func (db *DatabaseStore) GetMemberCount(ctx context.Context, isActive bool) (int, error) {
	var count int
	query := memberDbMethod.getMemberCount(isActive)
	err := db.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		log.Errorf("GetMemberCount query failed: %v", err)
		return 0, err
	}

	return count, nil
}

func (db *DatabaseStore) GetMembersPaginated(ctx context.Context, limit int, page int, active bool) ([]models.Member, error) {
	var members []models.Member
	query := memberDbMethod.getMembersPaginated(active)
	rows, err := db.pool.Query(ctx, query, limit, page*limit)
	if err != nil {
		log.Errorf("GetMembersPaginated query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	resourceMemo := make(map[string]models.MemberResource)

	for rows.Next() {
		var rIDs []string
		var member models.Member
		err = rows.Scan(&member.ID, &member.Name, &member.Email, &member.RFID, &member.Level, &rIDs, &member.SubscriptionID)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
			return nil, err
		}

		for _, rID := range rIDs {
			if resource, exist := resourceMemo[rID]; exist {
				member.Resources = append(member.Resources, resource)
				continue
			}

			resource, err := db.GetResourceByID(ctx, rID)
			if err != nil {
				logger.Errorf("error getting resource by id in memberResource lookup: %s %s\n", err.Error(), rID)
				continue
			}

			resourceMemo[rID] = models.MemberResource{
				ResourceID: resource.ID,
				Name:       resource.Name,
			}

			member.Resources = append(member.Resources, resourceMemo[rID])
		}

		members = append(members, member)
	}

	if rows.Err() != nil {
		log.Errorf("error iterating rows: %v", rows.Err())
		return nil, rows.Err()
	}

	return members, nil
}

func (db *DatabaseStore) GetMembers(ctx context.Context) []models.Member {
	var members []models.Member
	rows, err := db.pool.Query(ctx, memberDbMethod.getMember())
	if err != nil {
		log.Errorf("GetMembers failed: %v", err)
	}

	defer rows.Close()

	resourceMemo := make(map[string]models.MemberResource)

	for rows.Next() {
		var rIDs []string
		var member models.Member
		err = rows.Scan(&member.ID, &member.Name, &member.Email, &member.RFID, &member.Level, &rIDs, &member.SubscriptionID)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		// having issues with unmarshalling a jsonb object array from pgx
		// using a less efficient approach for now
		// TODO: fix this on the query level
		for _, rID := range rIDs {
			if _, exist := resourceMemo[rID]; exist {
				member.Resources = append(member.Resources, models.MemberResource{ResourceID: rID, Name: resourceMemo[rID].Name})
				continue
			}

			resource, err := db.GetResourceByID(ctx, rID)
			if err != nil {
				logger.Errorf("error getting resource by id in memberResource lookup: %s %s_\n", err.Error(), rID)
				continue
			}

			resourceMemo[rID] = models.MemberResource{
				ResourceID: resource.ID,
				Name:       resource.Name,
			}

			member.Resources = append(member.Resources, models.MemberResource{ResourceID: rID, Name: resource.Name})
		}

		members = append(members, member)
	}

	return members
}

func (db *DatabaseStore) getMemberBySubscriptionID(ctx context.Context, subscriptionID string) (models.Member, error) {
	for _, m := range db.GetMembers(ctx) {
		if m.SubscriptionID == subscriptionID {
			return m, nil
		}
	}

	return models.Member{}, fmt.Errorf("could not find member with subscriptionID: %s", subscriptionID)
}

// GetMemberByEmail - lookup a member by their email address
func (db *DatabaseStore) GetMemberByEmail(ctx context.Context, memberEmail string) (models.Member, error) {
	var member models.Member
	var rIDs []string

	err := db.pool.QueryRow(ctx, memberDbMethod.getMemberByEmail(), memberEmail).Scan(&member.ID, &member.Name, &member.Email, &member.RFID, &member.Level, &rIDs)
	if err == pgx.ErrNoRows {
		return member, err
	}
	if err != nil {
		log.Errorf("error getting member by email: %v", memberEmail)
		return member, fmt.Errorf("GetMemberByEmail failed: %w", err)
	}

	resourceMemo := make(map[string]models.MemberResource)

	// having issues with unmarshalling a jsonb object array from pgx
	// using a less efficient approach for now
	// TODO: fix this on the query level
	for _, rID := range rIDs {
		if _, exist := resourceMemo[rID]; exist {
			member.Resources = append(member.Resources, models.MemberResource{ResourceID: rID, Name: resourceMemo[rID].Name})
			continue
		}
		resource, err := db.GetResourceByID(ctx, rID)
		if err != nil {
			logger.Errorf("error getting resource by id in memberResource lookup: %s %s\n", err.Error(), rID)
		}

		resourceMemo[rID] = models.MemberResource{
			ResourceID: resource.ID,
			Name:       resource.Name,
		}
		member.Resources = append(member.Resources, models.MemberResource{ResourceID: rID, Name: resource.Name})
	}

	return member, nil
}

func (db *DatabaseStore) GetMemberByRFID(ctx context.Context, rfid string) (models.Member, error) {
	var member models.Member
	var rIDs []string

	err := db.pool.QueryRow(ctx, memberDbMethod.getMemberByRFID(), rfid).Scan(&member.ID, &member.Name, &member.Email, &member.RFID, &member.Level, &rIDs)
	if err == pgx.ErrNoRows {
		return member, err
	}
	if err != nil {
		log.Errorf("error getting member by email: %v", rfid)
		return member, fmt.Errorf("GetMemberByRFID failed: %w", err)
	}

	return member, nil
}

func (db *DatabaseStore) AssignRFID(ctx context.Context, email string, rfid string) (models.Member, error) {
	member, err := db.GetMemberByEmail(ctx, email)
	if err != nil {
		log.Errorf("error retrieving a member with that email address %s", err.Error())
		return member, err
	}

	err = db.pool.QueryRow(ctx, memberDbMethod.setMemberRFIDTag(), email, encodeRFID(rfid)).Scan(&member.RFID)
	if err != nil {
		return member, fmt.Errorf("AssignRFID failed: %v", err)
	}

	return member, err
}

func (db *DatabaseStore) UpdateMember(ctx context.Context, update models.Member) error {
	member, err := db.GetMemberByEmail(ctx, update.Email)
	if err != nil {
		log.Errorf("error retrieving a member with that email address %s", err.Error())
		return err
	}

	subID := member.SubscriptionID
	if len(update.SubscriptionID) > 0 {
		subID = update.SubscriptionID
	}

	commandTag, err := db.pool.Exec(ctx, memberDbMethod.updateMemberByEmail(), update.Name, subID, member.Email)
	if err != nil {
		return fmt.Errorf("UpdateMemberByEmail failed: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (db *DatabaseStore) UpdateMemberByID(ctx context.Context, memberID string, update models.Member) error {
	commandTag, err := db.pool.Exec(ctx, memberDbMethod.updateMemberByID(), update.Name, update.Email, update.SubscriptionID, memberID)
	if err != nil {
		return fmt.Errorf("UpdateMemberByID failed: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no member found with id: %s", memberID)
	}

	return nil
}

func (db *DatabaseStore) UpdateMemberBySubscriptionID(ctx context.Context, subscriptionID string, update models.Member) error {
	commandTag, err := db.pool.Exec(ctx, memberDbMethod.updateMemberBySubscriptionID(), update.Name, update.Email, subscriptionID)
	if err != nil {
		return fmt.Errorf("UpdateMemberBySubscriptionID failed: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no member found with subscriptionID: %s", subscriptionID)
	}

	return nil
}

func (db *DatabaseStore) AddNewMember(ctx context.Context, newMember models.Member) (models.Member, error) {
	// Give new members standard access at first; their actual status will be
	// evaluated on the next scheduled status check.
	if newMember.Level == 0 {
		newMember.Level = uint8(models.Standard)
	}

	commandTag, err := db.pool.Exec(ctx, memberDbMethod.insertMember(),
		newMember.Name, newMember.Email, newMember.Level, newMember.SubscriptionID)
	if err != nil {
		return models.Member{}, fmt.Errorf("AddNewMember query failed: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return models.Member{}, errors.New("no row affected")
	}

	if _, err := db.AddUserToDefaultResources(ctx, newMember.Email); err != nil {
		log.Error(err)
	}

	return newMember, nil
}

// GetMemberTiers - gets the member tiers from DB
func (db *DatabaseStore) GetTiers(ctx context.Context) []models.Tier {
	rows, err := db.pool.Query(ctx, tierDbMethod.getMemberTiers())
	if err != nil {
		log.Errorf("GetTiers failed: %v", err)
	}

	defer rows.Close()

	var tiers []models.Tier

	for rows.Next() {
		var t models.Tier
		err = rows.Scan(&t.ID, &t.Name)
		if err == nil {
			tiers = append(tiers, t)
		}
	}

	return tiers
}

var memberDbMethod MemberDatabaseMethod

// GetMembersWithCredit - gets members that have been credited a membership
//
//	if a member exists in the member_credits table
//	they are credited a membership
func (db *DatabaseStore) GetMembersWithCredit(ctx context.Context) []models.Member {
	rows, err := db.pool.Query(ctx, memberDbMethod.getMembersWithCredit())
	if err != nil {
		log.Errorf("error getting credited members: %v", err)
	}

	defer rows.Close()

	var members []models.Member

	for rows.Next() {
		var m models.Member
		err = rows.Scan(&m.ID, &m.Name, &m.Email, &m.RFID, &m.Level)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}

		members = append(members, m)
	}

	return members
}

// ProcessMember - add them member if they don't already exist.  Otherwise, make sure we have their name
func (db *DatabaseStore) ProcessMember(ctx context.Context, newMember models.Member) error {
	member, err := db.GetMemberByEmail(ctx, newMember.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, addErr := db.AddNewMember(ctx, newMember)
			return addErr
		}
		return err
	}

	if member.Name == "" {
		return db.updateMemberName(ctx, member.ID, newMember)
	}

	if member.SubscriptionID != newMember.SubscriptionID {
		return db.updateSubscriptionID(ctx, member.ID, newMember)
	}

	return nil
}

func (db *DatabaseStore) updateMemberName(ctx context.Context, memberID string, newMember models.Member) error {
	var member models.Member

	// if the member already exists, we might want to update their name.
	err := db.pool.QueryRow(ctx, memberDbMethod.updateMemberName(), memberID, newMember.Name).Scan(&member.Name)
	if err != nil {
		return fmt.Errorf("updateMemberName failed: %v", err)
	}

	return nil
}

func (db *DatabaseStore) updateSubscriptionID(ctx context.Context, memberID string, newMember models.Member) error {
	var member models.Member

	// if the member already exists, we might want to update their name.
	err := db.pool.QueryRow(ctx, memberDbMethod.updateMemberSubscriptionID(), memberID, newMember.SubscriptionID).Scan(&member.SubscriptionID)
	if err != nil {
		return fmt.Errorf("updateSubscriptionID failed: %v", err)
	}

	return nil
}

// SetMemberLevel sets a member's membership tier
func (db *DatabaseStore) SetMemberLevel(ctx context.Context, memberId string, level models.MemberLevel) error {
	rows, err := db.pool.Query(ctx, memberDbMethod.updateMembershipLevel(), memberId, level)
	if err != nil {
		log.Errorf("Set member level failed: %v", err)
		return err
	}
	defer rows.Close()
	return nil
}

// ApplyMemberCredits updates members tiers for all members with credit to Credited
func (db *DatabaseStore) ApplyMemberCredits(ctx context.Context) {
	//	Member credits are currently managed by DB commands.  #102 will address this.
	memberCredits := db.GetMembersWithCredit(ctx)
	for _, m := range memberCredits {
		err := db.SetMemberLevel(ctx, m.ID, models.Credited)
		if err != nil {
			log.Errorf("member credit failed: %v", err)
		}
	}
}

// UpdateMemberTiers updates member tiers based on the most recent payment amount
func (db *DatabaseStore) UpdateMemberTiers(ctx context.Context) {
	commandTag, err := db.pool.Exec(ctx, memberDbMethod.updateMemberTiers())
	if err != nil {
		log.Errorf("add members query failed: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		log.Errorf("no row affected")
	}
}

func (db *DatabaseStore) GetActiveMembersWithoutSubscription(ctx context.Context) []models.Member {
	var members []models.Member
	rows, err := db.pool.Query(ctx, memberDbMethod.getActiveMembersWithoutSubscription())
	if err != nil {
		log.Errorf("GetMembers failed: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var member models.Member
		err = rows.Scan(&member.ID, &member.Name, &member.Email, &member.RFID, &member.Level)
		if err != nil {
			log.Errorf("error scanning row: %s", err)
		}
		members = append(members, member)
	}
	return members
}
