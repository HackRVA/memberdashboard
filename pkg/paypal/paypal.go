package paypal

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"memberserver/internal/datastore"
	"memberserver/internal/models"
	"memberserver/internal/services/config"
	"memberserver/internal/services/logger"
	"net/http"
	"strings"
	"time"
)

type Paypal struct {
	db          datastore.DataStore
	accessToken string
}

type paypalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type Subscriber struct {
	Subscriber struct {
		ID        string `json:"id"`
		Summary   string `json:"summary"`
		EventType string `json:"event_type"`
		Name      struct {
			GivenName string `json:"given_name"`
			SurName   string `json:"surname"`
		} `json:"name"`
		Email string `json:"email_address"`
	} `json:"subscriber"`
}

type Payment struct {
	Amount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"amount"`
	Time time.Time `json:"time"`
}
type subscriptionResponse struct {
	Status      string `json:"status"`
	BillingInfo struct {
		LastPayment Payment `json:"last_payment"`
	} `json:"billing_info"`
}

func Setup(database datastore.DataStore) Paypal {
	return Paypal{
		db: database,
	}
}

// requestAccessToken - requests a BEARER access token to communicate with the api
func (p Paypal) requestAccessToken() (string, error) {
	var token string
	c, err := config.Load()

	if err != nil {
		return token, err
	}

	if len(c.PaypalClientID) == 0 {
		return token, fmt.Errorf("not a proper value for paypalClientID in the config")
	}
	if len(c.PaypalClientSecret) == 0 {
		return token, fmt.Errorf("not a proper value for paypalClientSecret in the config")
	}
	if len(c.PaypalURL) == 0 {
		return token, fmt.Errorf("not a proper value for paypalURL in the config")
	}

	payload := strings.NewReader("grant_type=client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.PaypalURL+"/v1/oauth2/token", payload)

	if err != nil {
		return token, err
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.PaypalClientID+":"+c.PaypalClientSecret)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()

	var newAccessToken paypalAccessTokenResponse

	err = json.NewDecoder(res.Body).Decode(&newAccessToken)
	if err != nil {
		return token, err
	}

	return newAccessToken.AccessToken, err
}

func (p Paypal) GetMemberFromSubscription(subscriptionID string) (models.Member, error) {
	var m models.Member
	var s Subscriber

	c, err := config.Load()
	if err != nil {
		fmt.Printf("error with config: %s", err)
		return m, err
	}
	url := fmt.Sprintf("%s/v1/billing/subscriptions/%s", c.PaypalURL, subscriptionID)
	token, err := p.requestAccessToken()
	if err != nil {
		logger.Errorf("error getting paypal access token %s\n", err.Error())
		return m, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return m, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return m, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		return m, err
	}

	logger.Debugf("member subscription response: %v", s)
	m.Email = s.Subscriber.Email
	m.Name = s.Subscriber.Name.GivenName + " " + s.Subscriber.Name.SurName
	m.SubscriptionID = subscriptionID

	return m, nil
}

func (p Paypal) GetSubscription(subscriptionID string) (string, models.Payment, error) {
	var lastPayment models.Payment
	c, err := config.Load()
	if err != nil {
		logger.Errorf("error with config: %s", err)
		return "", lastPayment, err
	}

	if len(c.PaypalURL) == 0 {
		return "", lastPayment, errors.New("no paypal url is set")
	}

	if len(p.accessToken) == 0 {
		p.accessToken, err = p.requestAccessToken()
		if err != nil {
			logger.Errorf("error getting paypal access token %s\n", err.Error())
			return "", lastPayment, err
		}
	}

	url := fmt.Sprintf("%s/v1/billing/subscriptions/%s", c.PaypalURL, subscriptionID)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", lastPayment, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", p.accessToken))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", lastPayment, err
	}
	defer res.Body.Close()

	var subscriptionStatus subscriptionResponse

	err = json.NewDecoder(res.Body).Decode(&subscriptionStatus)
	if err != nil {
		logger.Errorf("%s", err)
		return "", lastPayment, err
	}

	lastPayment.Amount = subscriptionStatus.BillingInfo.LastPayment.Amount.Value
	lastPayment.Time = subscriptionStatus.BillingInfo.LastPayment.Time

	return subscriptionStatus.Status, lastPayment, nil
}
