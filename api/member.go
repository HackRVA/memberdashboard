package api

import (
	"bytes"
	"encoding/json"
	"memberserver/api/models"
	"memberserver/datastore"
	"memberserver/resourcemanager"
	"memberserver/slack"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	log "github.com/sirupsen/logrus"
)

type MemberServer struct {
	store           datastore.DataStore
	ResourceManager *resourcemanager.ResourceManager
	AuthStrategy    union.Union
}

func (m *MemberServer) MemberEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		m.GetByEmailHandler(w, r)
	}

	if r.Method == http.MethodPut {
		m.UpdateMemberByEmailHandler(w, r)
	}
}

func (m *MemberServer) GetMembersHandler(w http.ResponseWriter, r *http.Request) {
	members := m.store.GetMembers()

	ok(w, members)
}

func (m *MemberServer) UpdateMemberByEmailHandler(w http.ResponseWriter, r *http.Request) {
	memberEmail := strings.TrimPrefix(r.URL.Path, "/api/member/email/")

	var request models.UpdateMemberRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if len(memberEmail) == 0 || !govalidator.IsEmail(memberEmail) {
		preconditionFailed(w, "invalid email")
		return
	}

	if err != nil {
		badRequest(w, err.Error())
		return
	}

	if len(request.FullName) == 0 {
		preconditionFailed(w, "fullName is required")
		return
	}

	_, err = m.store.GetMemberByEmail(memberEmail)

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	err = m.store.UpdateMemberByEmail(request.FullName, memberEmail)

	ok(w, models.EndpointSuccess{
		Ack: true,
	})

}

func (m *MemberServer) GetByEmailHandler(w http.ResponseWriter, r *http.Request) {
	memberEmail := strings.TrimPrefix(r.URL.Path, "/api/member/email/")

	if len(memberEmail) == 0 || !govalidator.IsEmail(memberEmail) {
		preconditionFailed(w, "invalid email")
		return
	}

	member, err := m.store.GetMemberByEmail(memberEmail)

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	ok(w, member)
}

func (m *MemberServer) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	_, user, _ := m.AuthStrategy.AuthenticateRequest(r)

	member, err := m.store.GetMemberByEmail(user.GetUserName())

	if err != nil {
		notFound(w, "error getting member by email")
		return
	}

	ok(w, member)
}

func (m *MemberServer) AssignRFIDHandler(w http.ResponseWriter, r *http.Request) {
	var assignRFIDRequest models.AssignRFIDRequest

	err := json.NewDecoder(r.Body).Decode(&assignRFIDRequest)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	m.assignRFID(w, assignRFIDRequest.Email, assignRFIDRequest.RFID)
}

func (m *MemberServer) AssignRFIDSelfHandler(w http.ResponseWriter, r *http.Request) {
	var assignRFIDRequest models.AssignRFIDRequest

	err := json.NewDecoder(r.Body).Decode(&assignRFIDRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, user, _ := m.AuthStrategy.AuthenticateRequest(r)

	m.assignRFID(w, user.GetUserName(), assignRFIDRequest.RFID)
}

func (m *MemberServer) GetTiersHandler(w http.ResponseWriter, r *http.Request) {
	tiers := m.store.GetTiers()

	ok(w, tiers)
}

func (m *MemberServer) assignRFID(w http.ResponseWriter, email, rfid string) {
	if len(rfid) == 0 {
		preconditionFailed(w, "not a valid rfid")
		return
	}

	m.removeMembersRFID(email)

	r, err := m.store.AssignRFID(email, rfid)
	if err != nil {
		notFound(w, "unable to assign rfid")
		return
	}

	ok(w, r)

	go m.ResourceManager.PushOne(models.Member{Email: email})
}

func (m *MemberServer) removeMembersRFID(email string) {
	member, err := m.store.GetMemberByEmail(email)
	if err != nil {
		log.Error(err)
		return
	}

	if member.RFID == "notset" || len(member.RFID) > 0 {
		return
	}

	for _, r := range member.Resources {
		resource, err := m.store.GetResourceByID(r.ResourceID)
		if err != nil {
			log.Error(err)
			continue
		}

		m.ResourceManager.RemoveMember(models.MemberAccess{
			Email:           member.Email,
			ResourceAddress: resource.Address,
			ResourceName:    resource.Name,
			Name:            member.Name,
			RFID:            member.RFID,
		})

	}
}

func (m *MemberServer) GetNonMembersOnSlackHandler(w http.ResponseWriter, r *http.Request) {
	nonMembers := slack.FindNonMembers(m.store)
	buf := bytes.NewBufferString(strings.Join(nonMembers[:], "\n"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=nonmembersOnSlack.csv")
	w.Write(buf.Bytes())
}

func (m *MemberServer) AddNewMemberHandler(w http.ResponseWriter, r *http.Request) {
	var newMember models.Member

	err := json.NewDecoder(r.Body).Decode(&newMember)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	addedMember, err := m.store.AddNewMember(newMember)
	if err != nil {
		http.Error(w, "error getting member by email", http.StatusNotFound)
		return
	}

	ok(w, addedMember)

	m.store.AssignRFID(addedMember.Email, addedMember.RFID)

	go m.ResourceManager.PushOne(addedMember)
}

func (m *MemberServer) GetMemberCounts(w http.ResponseWriter, req *http.Request) {
	chartType := req.URL.Query().Get("type")
	month := req.URL.Query().Get("month")

	if len(month) > 0 {
		date, err := time.Parse("", month)
		if err != nil {
			http.Error(w, "error looking up counts by month - use a valid date", http.StatusNotFound)
			return
		}
		ok(w, makeDistritutionChartByMonth(date, m.store))
		return
	}

	var charts []models.PaymentChart
	memberCounts, err := m.store.GetMemberCounts()
	if err != nil {
		http.Error(w, "error getting member counts", http.StatusNotFound)
		return
	}

	if len(chartType) > 0 {
		if chartType == "line" {
			charts = append(charts, makeMemberTrendChart(memberCounts))
		}

		if chartType == "pie" {
			charts = makeMemberDistributionChart(memberCounts)
		}
	}

	if len(chartType) == 0 {
		charts = append(charts, makeMemberTrendChart(memberCounts))
		charts = append(charts, makeMemberDistributionChart(memberCounts)...)
	}

	ok(w, charts)
}

func makeMemberTrendChart(counts []models.MemberCount) models.PaymentChart {
	var chart models.PaymentChart
	chart.Options.Title = "Membership Trends"
	chart.Type = "line"
	chart.Options.CurveType = "function"
	chart.Options.Legend = "bottom"
	chart.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "Member Count", Type: "number"}}

	for _, monthCount := range counts {
		var row []interface{}
		row = append(row, monthCount.Month.Format("Jan-06"))
		// explicitly exclude credited
		row = append(row, monthCount.Classic+monthCount.Standard+monthCount.Premium)
		chart.Rows = append(chart.Rows, row)
	}
	return chart
}

func makeDistritutionChartByMonth(month time.Time, store datastore.MemberStore) models.PaymentChart {
	var distributionChart models.PaymentChart
	memberCount, err := store.GetMemberCountByMonth(month)
	if err != nil {
		log.Errorf("error getting member counts")
		return distributionChart
	}

	var chart models.PaymentChart
	chart.Options.Title = month.Format("Jan-06")
	chart.Options.PieHole = 0.4
	chart.Type = "pie"

	chart.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "MemberLevelCount", Type: "number"}}

	levels := make(map[models.MemberLevel]uint8)
	levels[models.Credited] = uint8(memberCount.Credited)
	levels[models.Classic] = uint8(memberCount.Classic)
	levels[models.Standard] = uint8(memberCount.Standard)
	levels[models.Premium] = uint8(memberCount.Premium)

	for level, count := range levels {
		var row []interface{}
		row = append(row, (string)(models.MemberLevelToStr[level]))
		row = append(row, int(count))
		chart.Rows = append(chart.Rows, row)
	}

	return distributionChart
}

func makeMemberDistributionChart(counts []models.MemberCount) []models.PaymentChart {
	var distributionCharts []models.PaymentChart
	for _, monthCount := range counts {
		var chart models.PaymentChart
		chart.Options.Title = monthCount.Month.Format("Jan-06")
		chart.Options.PieHole = 0.4
		chart.Type = "pie"

		chart.Cols = []models.ChartCol{{Label: "Month", Type: "string"}, {Label: "MemberLevelCount", Type: "number"}}
		levels := make(map[models.MemberLevel]uint8)
		levels[models.Credited] = uint8(monthCount.Credited)
		levels[models.Classic] = uint8(monthCount.Classic)
		levels[models.Standard] = uint8(monthCount.Standard)
		levels[models.Premium] = uint8(monthCount.Premium)

		for level, count := range levels {
			var row []interface{}
			row = append(row, (string)(models.MemberLevelToStr[level]))
			row = append(row, int(count))
			chart.Rows = append(chart.Rows, row)
		}

		distributionCharts = append(distributionCharts, chart)
	}
	return distributionCharts
}
