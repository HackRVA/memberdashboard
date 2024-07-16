package services

import (
	"time"

	"github.com/HackRVA/memberserver/pkg/membermgr/models"
	"github.com/HackRVA/memberserver/pkg/membermgr/services/mail"
	"github.com/HackRVA/memberserver/pkg/mqtt"
	"github.com/sirupsen/logrus"

	go_mqtt "github.com/eclipse/paho.mqtt.golang"
)

type (
	Member interface {
		Add(models.Member) (models.Member, error)
		Get() []models.Member
		GetMembersWithLimit(limit int, offset int, active bool) []models.Member
		GetByEmail(email string) (models.Member, error)
		Update(models.Member) error
		AssignRFID(email string, rfid string) (models.Member, error)
		GetTiers() []models.Tier
		FindNonMembersOnSlack() []string
		GetMemberFromSubscription(subscriptionID string) (models.Member, error)
		CheckStatus(subscriptionID string) (models.Member, error)
		SetLevel(memberID string, level models.MemberLevel) error
		GetActiveMembersWithoutSubscription() []models.Member
	}

	MQTTHandler interface {
		// mqtt handlers
		HealthCheckHandler(client go_mqtt.Client, msg go_mqtt.Message)
		ReceiveHandler(client go_mqtt.Client, msg go_mqtt.Message)
		OnAccessEventHandler(payload models.LogMessage)
		OnHeartBeatHandler(client go_mqtt.Client, msg go_mqtt.Message)
		OnRemoveInvalidRequestHandler(client go_mqtt.Client, msg go_mqtt.Message)
	}

	Resource interface {
		MQTTHandler
		UpdateResourceACL(r models.Resource) error
		UpdateResources()
		EnableValidUIDs()
		RemovedInvalidUIDs()
		RemoveMember(memberAccess models.MemberAccess)
		Open(resource models.Resource)
		RemoveOne(member models.Member)
		PushOne(m models.Member)
		DeleteResourceACL()
		CheckStatus(r models.Resource)
		MQTT() mqtt.MQTTServer
	}

	Logger interface {
		SetLevel(level logrus.Level)
		Println(args ...interface{})
		Printf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Fatalf(format string, args ...interface{})
		Tracef(format string, args ...interface{})
		Print(args ...interface{})
		Error(args ...interface{})
		Debug(args ...interface{})
		Info(args ...interface{})
		Fatal(args ...interface{})
		Trace(args ...interface{})
	}

	Mailer interface {
		SendCommunication(communication mail.CommunicationTemplate, recipient string, model interface{}) (bool, error)
		IsThrottled(c models.Communication, member models.Member) bool
	}

	Report interface {
		GetAccessStatsChart(date time.Time, resourceName string) (models.ReportChart, error)
		GetMemberChurn() (int, error)
		GetMemberCountsChartByMonth(date time.Time) models.ReportChart
		GetMemberCountsCharts(chartType string) ([]models.ReportChart, error)
	}

	Job interface {
		CheckActiveMembersWithoutSubscription()
		CheckMemberSubscriptions()
		CheckResourceInit()
		CheckResourceInterval()
		CheckIPAddressInterval()
		RemovedInvalidUIDs()
		EnableValidUIDs()
		UpdateResources()
		UpdateMemberCounts()
	}

	Scheduler interface {
		Setup(j Job)
		scheduleTask(interval time.Duration, initFunc func(), tickFunc func())
	}
)
