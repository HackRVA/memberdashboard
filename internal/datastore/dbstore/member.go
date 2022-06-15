package dbstore

import (
	"context"
	"errors"
	"fmt"
	"memberserver/internal/models"
	"memberserver/internal/services/logger"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func (db *DatabaseStore) GetMembers() []models.Member {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var members []models.Member
	rows, err := dbPool.Query(db.ctx, memberDbMethod.getMember())
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

			resource, err := db.GetResourceByID(rID)
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

// GetMemberByEmail - lookup a member by their email address
func (db *DatabaseStore) GetMemberByEmail(memberEmail string) (models.Member, error) {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var member models.Member
	var rIDs []string

	err = dbPool.QueryRow(context.Background(), memberDbMethod.getMemberByEmail(), memberEmail).Scan(&member.ID, &member.Name, &member.Email, &member.RFID, &member.Level, &rIDs)
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
		resource, err := db.GetResourceByID(rID)
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

func (db *DatabaseStore) GetMemberByRFID(rfid string) (models.Member, error) {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var member models.Member
	var rIDs []string

	err = dbPool.QueryRow(context.Background(), memberDbMethod.getMemberByRFID(), rfid).Scan(&member.ID, &member.Name, &member.Email, &member.RFID, &member.Level, &rIDs)
	if err == pgx.ErrNoRows {
		return member, err
	}
	if err != nil {
		log.Errorf("error getting member by email: %v", rfid)
		return member, fmt.Errorf("GetMemberByRFID failed: %w", err)
	}

	return member, nil
}

func (db *DatabaseStore) AssignRFID(email string, rfid string) (models.Member, error) {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	member, err := db.GetMemberByEmail(email)
	if err != nil {
		log.Errorf("error retrieving a member with that email address %s", err.Error())
		return member, err
	}

	err = dbPool.QueryRow(context.Background(), memberDbMethod.setMemberRFIDTag(), email, encodeRFID(rfid)).Scan(&member.RFID)
	if err != nil {
		return member, fmt.Errorf("AssignRFID failed: %v", err)
	}

	return member, err
}

func (db *DatabaseStore) UpdateMember(update models.Member) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	member, err := db.GetMemberByEmail(update.Email)

	if err != nil {
		log.Errorf("error retrieving a member with that email address %s", err.Error())
		return err
	}

	subID := member.SubscriptionID
	if len(update.SubscriptionID) > 0 {
		subID = update.SubscriptionID
	}

	commandTag, err := dbPool.Exec(context.Background(), memberDbMethod.updateMemberByEmail(), update.Name, subID, member.Email)
	if err != nil {
		return fmt.Errorf("UpdateMemberByEmail failed: %v", err)
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (db *DatabaseStore) AddNewMember(newMember models.Member) (models.Member, error) {
	err := db.AddMembers([]models.Member{newMember})
	if err != nil {
		return models.Member{}, err
	}
	return newMember, nil
}

// GetMemberTiers - gets the member tiers from DB
func (db *DatabaseStore) GetTiers() []models.Tier {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	rows, err := dbPool.Query(context.Background(), tierDbMethod.getMemberTiers())
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
//  if a member exists in the member_credits table
//  they are credited a membership
func (db *DatabaseStore) GetMembersWithCredit() []models.Member {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	rows, err := dbPool.Query(db.ctx, memberDbMethod.getMembersWithCredit())
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

// AddMembers adds multiple members to the DatabaseStore
func (db *DatabaseStore) AddMembers(members []models.Member) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	sqlStr := `INSERT INTO membership.members(
name, email, member_tier_id, subscription_id)
VALUES `

	var valStr []string
	for _, m := range members {
		// postgres doesn't like apostrophes
		memberName := strings.Replace(m.Name, "'", "''", -1)

		// Give new members standard access at first
		//   Their actual status will be evaluated the next day
		if m.Level == 0 {
			m.Level = uint8(models.Standard)
		}

		valStr = append(valStr, fmt.Sprintf("('%s', '%s', %d, '%s')", memberName, m.Email, m.Level, m.SubscriptionID))
	}

	str := strings.Join(valStr, ",")

	commandTag, err := dbPool.Exec(context.Background(), sqlStr+str+"ON CONFLICT DO NOTHING;")
	if err != nil {
		return fmt.Errorf("add members query failed: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	for _, m := range members {
		log.Info("Adding default resource")
		db.AddUserToDefaultResources(m.Email)
	}

	return err
}

// ProcessMember - add them member if they don't already exist.  Otherwise, make sure we have their name
func (db *DatabaseStore) ProcessMember(newMember models.Member) error {
	member, err := db.GetMemberByEmail(newMember.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return db.AddMembers([]models.Member{newMember})
		}
		return err
	}

	if member.Name == "" {
		return db.updateMemberName(member.ID, newMember)
	}

	if member.SubscriptionID != newMember.SubscriptionID {
		return db.updateSubscriptionID(member.ID, newMember)
	}

	return nil
}

func (db *DatabaseStore) updateMemberName(memberID string, newMember models.Member) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var member models.Member

	// if the member already exists, we might want to update their name.
	err = dbPool.QueryRow(context.Background(), memberDbMethod.updateMemberName(), memberID, newMember.Name).Scan(&member.Name)
	if err != nil {
		return fmt.Errorf("updateMemberName failed: %v", err)
	}

	return nil
}

func (db *DatabaseStore) updateSubscriptionID(memberID string, newMember models.Member) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	var member models.Member

	// if the member already exists, we might want to update their name.
	err = dbPool.QueryRow(context.Background(), memberDbMethod.updateMemberSubscriptionID(), memberID, newMember.SubscriptionID).Scan(&member.SubscriptionID)
	if err != nil {
		return fmt.Errorf("updateSubscriptionID failed: %v", err)
	}

	return nil
}

// SetMemberLevel sets a member's membership tier
func (db *DatabaseStore) SetMemberLevel(memberId string, level models.MemberLevel) error {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	rows, err := dbPool.Query(context.Background(), memberDbMethod.updateMembershipLevel(), memberId, level)
	if err != nil {
		log.Errorf("Set member level failed: %v", err)
		return err
	}
	defer rows.Close()
	return nil
}

// ApplyMemberCredits updates members tiers for all members with credit to Credited
func (db *DatabaseStore) ApplyMemberCredits() {
	//	Member credits are currently managed by DB commands.  #102 will address this.
	memberCredits := db.GetMembersWithCredit()
	for _, m := range memberCredits {
		err := db.SetMemberLevel(m.ID, models.Credited)
		if err != nil {
			log.Errorf("member credit failed: %v", err)
		}
	}
}

// UpdateMemberTiers updates member tiers based on the most recent payment amount
func (db *DatabaseStore) UpdateMemberTiers() {
	dbPool, err := pgxpool.Connect(db.ctx, db.connectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

	commandTag, err := dbPool.Exec(context.Background(), memberDbMethod.updateMemberTiers())
	if err != nil {
		log.Errorf("add members query failed: %v", err)
	}
	if commandTag.RowsAffected() == 0 {
		log.Errorf("no row affected")
	}

}
