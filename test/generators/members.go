package generators

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/HackRVA/memberserver/datastore"
	"github.com/HackRVA/memberserver/models"

	log "github.com/sirupsen/logrus"
	"syreclabs.com/go/faker"
)

func Seed(db datastore.DataStore, numMembers int) {
	ctx := context.Background()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	FakeResources(db)

	if _, err := db.AddNewMember(ctx, TestMember()); err != nil {
		log.Errorf("error adding test members: %s", err)
	}

	for i := 0; i < numMembers; i++ {
		member := FakeMember(rng)
		if _, err := db.AddNewMember(ctx, member); err != nil {
			log.Errorf("error adding test members: %s", err)
		}
		log.Printf("Added member %v", member.Name)
		if member.Level > 1 {
			member, _ = db.GetMemberByEmail(ctx, member.Email)
			memberLevelID, _ := strconv.Atoi(faker.Number().Between(1, 5))
			if err := db.SetMemberLevel(ctx, member.ID, models.MemberLevel(memberLevelID)); err != nil {
				log.Errorf("error setting member level: %s", err)
			}
		}
	}

	FakeMemberCounts(24, db)
	FakeAccessEvents(50, db)
	RegisterTestUser(db)
}

func FakeAccessEvents(numOfEvents int, db datastore.DataStore) {
	ctx := context.Background()
	resources := db.GetResources(ctx)

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
			if err := db.LogAccessEvent(ctx, logMsg); err != nil {
				log.Errorf("error logging event: %s", err)
			}
			log.Infof("Added log event for %s time: %s", logMsg.Username, eventTime)
		}
	}
}

// FakeMember generates a random member using a local RNG instance
func FakeMember(rng *rand.Rand) models.Member {
	level := uint8(rng.Intn(5) + 1) // Level between 1 and 5
	resources := []models.MemberResource{}
	return models.Member{
		Name:           faker.Name().Name(),
		Email:          faker.Internet().Email(),
		Level:          level,
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
	if err := db.RegisterUser(context.Background(), models.Credentials{
		Email:    "test@test.com",
		Password: "test",
	}); err != nil {
		log.Error(err)
	}
}

func FakeResources(db datastore.DataStore) {
	ctx := context.Background()
	if _, err := db.RegisterResource(ctx, faker.App().Name(), string(faker.Internet().IpV4Address()), false); err != nil {
		log.Error(err)
	}
	if _, err := db.RegisterResource(ctx, faker.App().Name(), string(faker.Internet().IpV4Address()), true); err != nil {
		log.Error(err)
	}
}

func FakeMemberCounts(numberOfMonths int, db datastore.DataStore) {
	// var months []models.MemberCount
	//
	// for i := 1; i < numberOfMonths; i++ {
	// 	m := time.Now().AddDate(0, -i, 0)
	// 	months = append(months, models.MemberCount{
	// 		Month:    m,
	// 		Classic:  faker.Number().NumberInt(3),
	// 		Standard: faker.Number().NumberInt(3),
	// 		Premium:  faker.Number().NumberInt(3),
	// 		Credited: faker.Number().NumberInt(3),
	// 	})
	// }
}
