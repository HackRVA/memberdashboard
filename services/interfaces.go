package services

import (
	"context"
	"time"

	"github.com/HackRVA/memberserver/models"
	"github.com/HackRVA/memberserver/pkg/mqtt"
	"github.com/HackRVA/memberserver/services/mail"
	"github.com/sirupsen/logrus"

	go_mqtt "github.com/eclipse/paho.mqtt.golang"
)

type (
	Member interface {
		Add(ctx context.Context, m models.Member) (models.Member, error)
		Get(ctx context.Context) []models.Member
		GetMembersPaginated(ctx context.Context, limit int, offset int, active bool) []models.Member
		GetMemberCount(ctx context.Context, isActive bool) (int, error)
		GetByEmail(ctx context.Context, email string) (models.Member, error)
		Update(ctx context.Context, m models.Member) error
		UpdateMemberByID(ctx context.Context, memberID string, update models.Member) error
		AssignRFID(ctx context.Context, email string, rfid string) (models.Member, error)
		GetTiers(ctx context.Context) []models.Tier
		FindNonMembersOnSlack(ctx context.Context) []string
		GetMemberFromSubscription(subscriptionID string) (models.Member, error)
		CheckStatus(ctx context.Context, subscriptionID string) (models.Member, error)
		SetLevel(ctx context.Context, memberID string, level models.MemberLevel) error
		GetActiveMembersWithoutSubscription(ctx context.Context) []models.Member
	}

	MQTTHandler interface {
		// mqtt handlers
		HealthCheckHandler(client go_mqtt.Client, msg go_mqtt.Message)
		ReceiveHandler(client go_mqtt.Client, msg go_mqtt.Message)
		OnAccessEventHandler(payload models.LogMessage)
		OnHeartBeatHandler(client go_mqtt.Client, msg go_mqtt.Message)
		OnRemoveInvalidRequestHandler(client go_mqtt.Client, msg go_mqtt.Message)
	}

	ResourceUpdater interface {
		PushOne(m models.Member)
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
		SendCommunication(ctx context.Context, communication mail.CommunicationTemplate, recipient string, model interface{}) (bool, error)
		IsThrottled(ctx context.Context, c models.Communication, member models.Member) bool
	}

	Report interface {
		GetAccessStatsChart(ctx context.Context, date time.Time, resourceName string) (models.ReportChart, error)
		GetMemberChurn(ctx context.Context) (int, error)
		GetMemberCountsChartByMonth(ctx context.Context, date time.Time) models.ReportChart
		GetMemberCountsCharts(ctx context.Context, chartType string) ([]models.ReportChart, error)
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
