package paypal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type Paypal struct {
	accessToken string
	logger      Logger
	config      Config
}

type Config struct {
	url      string
	clientID string
	secret   string
}

type paypalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
type subscriptionResponse struct {
	Status      string `json:"status"`
	BillingInfo struct {
		LastPayment payment `json:"last_payment"`
	} `json:"billing_info"`
	Subscriber Subscriber `json:"subscriber"`
}

type Subscriber struct {
	ID        string `json:"id"`
	Summary   string `json:"summary"`
	EventType string `json:"event_type"`
	Name      struct {
		GivenName string `json:"given_name"`
		SurName   string `json:"surname"`
	} `json:"name"`
	Email string `json:"email_address"`
}

type payment struct {
	Amount struct {
		CurrencyCode string `json:"currency_code"`
		Value        string `json:"value"`
	} `json:"amount"`
	Time time.Time `json:"time"`
}

type Payment struct {
	Amount string    `json:"amount"`
	Time   time.Time `json:"time"`
}

func Setup(url string, clientID string, clientSecret string, logger Logger) Paypal {
	if logger == nil {
		// allows for custom loggers to be passed in
		logger = logrus.New()
	}
	return Paypal{
		logger: logger,
		config: Config{
			url:      url,
			clientID: clientID,
			secret:   clientSecret,
		},
	}
}

// requestAccessToken - requests a BEARER access token to communicate with the api
func (p Paypal) requestAccessToken() (string, error) {
	var token string
	if err := p.checkConfig(); err != nil {
		return "", err
	}

	payload := strings.NewReader("grant_type=client_credentials")

	client := &http.Client{}
	req, err := http.NewRequest("POST", p.config.url+"/v1/oauth2/token", payload)

	if err != nil {
		return token, err
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(p.config.clientID+":"+p.config.secret)))
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

func (p Paypal) GetSubscription(subscriptionID string) (status string, lastPaymentAmount string, lastPaymentTime time.Time, err error) {
	s, err := p.getSubscription(subscriptionID)
	if err != nil {
		return status, lastPaymentAmount, lastPaymentTime, err
	}

	lastPaymentAmount = s.BillingInfo.LastPayment.Amount.Value
	lastPaymentTime = s.BillingInfo.LastPayment.Time

	return status, lastPaymentAmount, lastPaymentTime, nil
}

func (p Paypal) GetSubscriber(subscriptionID string) (name string, email string, err error) {
	s, err := p.getSubscription(subscriptionID)
	if err != nil {
		return name, email, err
	}

	name = s.Subscriber.Name.GivenName + " " + s.Subscriber.Name.SurName
	email = s.Subscriber.Email

	return name, email, nil
}

func (p Paypal) getSubscription(subscriptionID string) (response subscriptionResponse, err error) {
	if len(p.accessToken) == 0 {
		p.accessToken, err = p.requestAccessToken()
		if err != nil {
			p.logger.Errorf("error getting paypal access token %s\n", err.Error())
			return response, err
		}
	}
	url := fmt.Sprintf("%s/v1/billing/subscriptions/%s", p.config.url, subscriptionID)
	token, err := p.requestAccessToken()
	if err != nil {
		p.logger.Errorf("error getting paypal access token %s\n", err.Error())
		return response, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return response, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&response)
	return response, err
}

func (p Paypal) checkConfig() error {
	if len(p.config.clientID) == 0 {
		return fmt.Errorf("not a proper value for paypalClientID in the config")
	}
	if len(p.config.secret) == 0 {
		return fmt.Errorf("not a proper value for paypalClientSecret in the config")
	}
	if len(p.config.url) == 0 {
		return fmt.Errorf("not a proper value for paypalURL in the config")
	}

	return nil
}
