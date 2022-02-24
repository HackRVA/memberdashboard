package main

import (
	"context"
	"fmt"
	"math/rand"
	"memberserver/api/models"
	"memberserver/config"
	"memberserver/datastore/dbstore"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"

	"github.com/Rhymond/go-money"
	"syreclabs.com/go/faker"
)

var tiers = []int64{0, 0, 0, 30, 35, 50}

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

	for i := 0; i < count; i++ {
		member := fakeMember()
		db.AddMembers([]models.Member{member})
		log.Printf("Added member %v", member.Name)
		if member.Level > 1 {
			member, _ = db.GetMemberByEmail(member.Email)
			lastPayment := time.Now().AddDate(0, 0, -rand.Intn(70))
			numPayments := rand.Intn(6)
			log.Printf("Creating %v payments", numPayments)
			payments := fakePaymentHistory(member, lastPayment, numPayments)
			db.AddPayments(payments)
		}
	}

	fakeMemberCounts(24)
	fakeAccessEvents(50)
	registerTestUser()
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

func fakePaymentHistory(member models.Member, lastPayment time.Time, numberOfPayments int) []models.Payment {
	payments := []models.Payment{}
	for i := 0; i < numberOfPayments; i++ {
		paymentDate := lastPayment.AddDate(0, -i, 0)
		payments = append(payments, models.Payment{
			ID:       faker.Number().Number(8),
			Date:     paymentDate,
			Amount:   *money.New(tiers[member.Level]*100, "USD"),
			Provider: 1,
			MemberID: member.ID,
			Email:    member.Email,
			Name:     member.Name,
		})
	}
	return payments
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

	_, err = dbPool.Exec(context.Background(), sqlStr+str+" ON CONFLICT DO NOTHING;")
	if err != nil {
		log.Errorf("conn.Exec failed: %v", err)
	}
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
