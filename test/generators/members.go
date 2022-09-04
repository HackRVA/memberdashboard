package main

import (
	"context"
	"fmt"
	"math/rand"
	"memberserver/internal/datastore/dbstore"
	"memberserver/internal/models"
	"memberserver/internal/services/config"
	"memberserver/internal/services/scheduler/jobs"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"syreclabs.com/go/faker"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Provide an argument for # of members to create.")
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to parse %v as count.", os.Args[1])
	}
	fakeResources()

	rand.Seed(time.Now().UnixNano())
	db, _ := dbstore.Setup()
	db.AddMembers([]models.Member{testMember()})

	jobManager := jobs.New(db, log.New())

	for i := 0; i < count; i++ {
		member := fakeMember()
		db.AddMembers([]models.Member{member})
		log.Printf("Added member %v", member.Name)
		if member.Level > 1 {
			member, _ = db.GetMemberByEmail(member.Email)
			memberLevelID, _ := strconv.Atoi(faker.Number().Between(1, 5))
			jobManager.SetMemberLevel(models.ActiveStatus, models.Payment{
				Amount: strconv.Itoa(int(models.MemberLevelToAmount[models.MemberLevel(memberLevelID)])),
				Time:   time.Now().AddDate(0, 0, -rand.Intn(70)),
			}, member)
		}
	}

	fakeMemberCounts(24)
	fakeAccessEvents(50)
	registerTestUser()
}

func fakeAccessEvents(numOfEvents int) {
	db, _ := dbstore.Setup()
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

func fakeMember() models.Member {
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

func testMember() models.Member {
	return models.Member{
		Name:           "test",
		Email:          "test@test.com",
		Level:          uint8(models.Premium),
		RFID:           faker.Lorem().Characters(10),
		Resources:      []models.MemberResource{},
		SubscriptionID: faker.Internet().MacAddress(),
	}
}

func registerTestUser() {
	db, _ := dbstore.Setup()
	db.RegisterUser(models.Credentials{
		Email:    "test@test.com",
		Password: "test",
	})
}

func fakeResources() {
	db, _ := dbstore.Setup()
	db.RegisterResource(faker.App().Name(), string(faker.Internet().IpV4Address()), false)
	db.RegisterResource(faker.App().Name(), string(faker.Internet().IpV4Address()), true)
}

func fakeMemberCounts(numberOfMonths int) {
	conf, _ := config.Load()

	dbPool, err := pgxpool.Connect(context.Background(), conf.DBConnectionString)
	if err != nil {
		log.Printf("got error: %v\n", err)
	}
	defer dbPool.Close()

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

	var valStr []string

	sqlStr := `INSERT INTO membership.member_counts(month, classic, standard, premium, credited)
VALUES `

	for _, p := range months {
		valStr = append(valStr, fmt.Sprintf("(TO_DATE('%s', 'YYYYMM'), %d, %d, %d, %d)", p.Month.Format("200601"), p.Classic, p.Standard, p.Premium, p.Credited))
	}

	str := strings.Join(valStr, ",")

	log.Infof("Adding %d months of member counts", len(months))

	commandTag, err := dbPool.Exec(context.Background(), sqlStr+str+" ON CONFLICT DO NOTHING;")
	if err != nil {
		log.Errorf("conn.Exec failed: %v", err)
	}
	if commandTag.RowsAffected() != 1 {
		log.Errorf("no row affected")
	}
}
