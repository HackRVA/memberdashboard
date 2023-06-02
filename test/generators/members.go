package generators

import (
	"math/rand"

	"github.com/HackRVA/memberserver/pkg/membermgr/datastore"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"

	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"syreclabs.com/go/faker"
)

func Seed(db datastore.DataStore, numMembers int) {

	FakeResources(db)

	rand.Seed(time.Now().UnixNano())
	db.AddMembers([]models.Member{TestMember()})

	for i := 0; i < numMembers; i++ {
		member := FakeMember()
		db.AddMembers([]models.Member{member})
		log.Printf("Added member %v", member.Name)
		if member.Level > 1 {
			member, _ = db.GetMemberByEmail(member.Email)
			memberLevelID, _ := strconv.Atoi(faker.Number().Between(1, 5))
			db.SetMemberLevel(member.ID, models.MemberLevel(memberLevelID))
		}
	}

	FakeMemberCounts(24, db)
	FakeAccessEvents(50, db)
	RegisterTestUser(db)
}

func FakeAccessEvents(numOfEvents int, db datastore.DataStore) {
	resources := db.GetResources()

	for resourceIndex, r := range resources {
		if resourceIndex == 5 {
			break
		}
		for i := 0; i < numOfEvents; i++ {
			eventTime := faker.Time().Between(time.Now().Add((time.Hour*24)*-30), time.Now())
			logMsg := models.LogMessage{
				Type:      "access",
				EventTime: eventTime.Unix(),
				IsKnown:   "true",
				Username:  faker.Name().Name(),
				RFID:      string(faker.Internet().IpV4Address()),
				Door:      r.Name,
			}
			err := db.LogAccessEvent(logMsg)
			if err != nil {
				log.Errorf("error logging event: %s", err)
			}
			log.Infof("Added log event for %s time: %s", logMsg.Username, eventTime)
		}
	}
}

func FakeMember() models.Member {
	level, _ := strconv.Atoi(faker.Number().Between(1, 5))
	resources := []models.MemberResource{}
	return models.Member{
		Name:           faker.Name().Name(),
		Email:          faker.Internet().Email(),
		Level:          uint8(level),
		RFID:           faker.Lorem().Characters(10),
		Resources:      resources,
		SubscriptionID: faker.Internet().MacAddress(),
	}
}

func TestMember() models.Member {
	return models.Member{
		Name:  "test",
		Email: "test@test.com",
		Level: uint8(models.Premium),
		RFID:  faker.Lorem().Characters(10),
		Resources: []models.MemberResource{
			{
				Name: "admin",
			},
		},
		SubscriptionID: faker.Internet().MacAddress(),
	}
}

func RegisterTestUser(db datastore.DataStore) {
	db.RegisterUser(models.Credentials{
		Email:    "test@test.com",
		Password: "test",
	})
}

func FakeResources(db datastore.DataStore) {
	db.RegisterResource(faker.App().Name(), string(faker.Internet().IpV4Address()), false)
	db.RegisterResource(faker.App().Name(), string(faker.Internet().IpV4Address()), true)
}

func FakeMemberCounts(numberOfMonths int, db datastore.DataStore) {
	var months []models.MemberCount

	for i := 1; i < numberOfMonths; i++ {
		m := time.Now().AddDate(0, -i, 0)
		months = append(months, models.MemberCount{
			Month:    m,
			Classic:  faker.Number().NumberInt(3),
			Standard: faker.Number().NumberInt(3),
			Premium:  faker.Number().NumberInt(3),
			Credited: faker.Number().NumberInt(3),
		})
	}
}
