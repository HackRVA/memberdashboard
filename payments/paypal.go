package payments

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"memberserver/api/models"
	"memberserver/config"

	"github.com/Rhymond/go-money"
	log "github.com/sirupsen/logrus"
	"gopkg.in/errgo.v2/fmt/errors"
)

type paypalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type paypalTransactions struct {
	Transactions []paypalTransaction `json:"transaction_details"`
}

type paypalTransaction struct {
	Transaction transaction `json:"transaction_info"`
	Payer       payer       `json:"payer_info"`
}

type transaction struct {
	Status  string            `json:"transaction_status"`
	Subject string            `json:"transaction_subject"`
	Date    string            `json:"transaction_initiation_date"`
	Amount  transactionAmount `json:"transaction_amount"`
}

type payer struct {
	Email string    `json:"email_address"`
	Name  payerName `json:"payer_name"`
}

type payerName struct {
	GivenName string `json:"given_name"`
	Surname   string `json:"surname"`
	FullName  string `json:"alternate_full_name"`
}

type transactionAmount struct {
	CurrencyCode string  `json:"currency_code"`
	Value        float64 `json:"value,string"`
}

var accessToken = ""

func getPaypalPayments(startDate string, endDate string) ([]models.Payment, error) {
	var payments []models.Payment
	c, err := config.Load()
	if err != nil {
		fmt.Printf("error with config: %s", err)
		return payments, err
	}

	url := fmt.Sprintf("%s/v1/reporting/transactions?start_date=%s&end_date=%s&fields=transaction_info,payer_info", c.PaypalURL, startDate, endDate)
	token, err := requestPaypalAccessToken()
	if err != nil {
		log.Errorf("error getting paypal access token %s\n", err.Error())
		return payments, err
	}

	if len(token) == 0 {
		return payments, errors.Newf("invalid token from paypal: %s", token)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return payments, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return payments, err
	}
	defer res.Body.Close()

	var transactions paypalTransactions

	err = json.NewDecoder(res.Body).Decode(&transactions)
	if err != nil {
		return payments, err
	}

	timeLayout := "2006-01-02T15:04:05-0700"

	for _, t := range transactions.Transactions {
		if t.Transaction.Status != "S" {
			// found a transaction, but it didn't have a Status Code of 'S'
			//
			//   From PayPal:
			// Status code	Description
			// 	D	PayPal or merchant rules denied the transaction.
			// 	F	The original recipient partially refunded the transaction.
			// 	P	The transaction is pending. The transaction was created but waits for another payment process to complete, such as an ACH transaction, before the status changes to S.
			// 	S	The transaction successfully completed without a denial and after any pending statuses.
			// 	V	A successful transaction was reversed and funds were refunded to the original sender.
			continue
		}

		var p models.Payment

		// i don't think this will handle pennies, but i don't think it matters
		p.Amount = *money.New((int64)(t.Transaction.Amount.Value)*100, t.Transaction.Amount.CurrencyCode)
		p.Email = t.Payer.Email
		p.Name = t.Payer.Name.FullName
		p.Date, err = time.Parse(timeLayout, t.Transaction.Date)
		if err != nil {
			log.Errorf("error in date of a transaction: %s\n", err.Error())
		}
		p.Provider = models.Paypal

		payments = append(payments, p)
	}

	return payments, err
}

// requestPaypalAccessToken - requests a BEARER access token to communicate with the api
func requestPaypalAccessToken() (string, error) {
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

	token = newAccessToken.AccessToken

	return newAccessToken.AccessToken, err
}
