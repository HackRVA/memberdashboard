package paymentproviders

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dfirebaugh/memberserver/config"
)

type paypalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getPaypalPayments(startDate time.Time, endDate time.Time) {
	c, err := config.Load(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))

	if err != nil {
		fmt.Printf("error with config: %s", err)
	}

	token, err := requestPaypalAccessToken()
	if err != nil {
		fmt.Printf("error getting access token %s\n", err)
	}

	u, err := url.Parse(c.PaypalURL + "/v1/reporting/transactions")
	if err != nil {
		fmt.Printf("error with paypalURL: %s", err)
	}

	u.Scheme = "https"
	q := u.Query()
	q.Set("start_date", startDate.String())
	q.Set("end_date", endDate.String())

	u.RawQuery = q.Encode()

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, u.RequestURI(), nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// requestPaypalAccessToken - requests a BEARER access token to communicate with the api
func requestPaypalAccessToken() (string, error) {
	var token string
	c, err := config.Load(os.Getenv("MEMBER_SERVER_CONFIG_FILE"))

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
		fmt.Println(err)
		return token, err
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.PaypalClientID+" "+c.PaypalClientSecret)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
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
